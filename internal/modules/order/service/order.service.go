package order_service

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/gateway"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	loyalty_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty/services"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_filters "larsa-e-commerce-microservice-v2/internal/modules/order/filters"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_repos "larsa-e-commerce-microservice-v2/internal/modules/order/repos"
	order_utils "larsa-e-commerce-microservice-v2/internal/modules/order/utils"
	settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"
	settings_sales_service "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: things needed to be done
// 1. apply assign to driver feature with status
// 2. apply confirm delivery feature with status
type OrderService interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*order_models.OrderResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string, limit int64) ([]order_models.OrderResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*types.DataWithPagination[order_models.OrderResponse], *helpers.ApiError)

	Create(ctx context.Context, dto *order_dto.CreateOrderDTO) (*order_models.OrderResponse, *helpers.ApiError)

	// Actions
	UpdateOrderStatus(ctx context.Context, id primitive.ObjectID, newStatusName string) (*order_models.OrderResponse, *helpers.ApiError)
	UpdateAdminSeen(ctx context.Context, id primitive.ObjectID, dto *order_dto.UpdateAdminSeenDTO) (*order_models.OrderResponse, *helpers.ApiError)
	UpdateTransitionStatus(ctx context.Context, id primitive.ObjectID, status string) (*order_models.OrderResponse, *helpers.ApiError)

	// Admin Actions
	DeleteAllOrders(ctx context.Context) *helpers.ApiError
}

type orderService struct {
	repo                     order_repos.OrderRepo
	settingsSalesService     settings_sales_service.SalesSettingsService
	orderCancelRefundService OrderCancelRefundService
	loyaltyOrderService      loyalty_services.LoyaltyOrderService
	rolesGw                  *gateway.RolesGw
}

func NewOrderService(i *do.Injector) (OrderService, error) {
	repo := do.MustInvoke[order_repos.OrderRepo](i)
	settingsSalesService := do.MustInvoke[settings_sales_service.SalesSettingsService](i)
	orderCancelRefundService := do.MustInvoke[OrderCancelRefundService](i)
	loyaltyOrderService := do.MustInvoke[loyalty_services.LoyaltyOrderService](i)
	rolesGw := do.MustInvoke[*gateway.RolesGw](i)
	return &orderService{
		repo:                     repo,
		settingsSalesService:     settingsSalesService,
		orderCancelRefundService: orderCancelRefundService,
		loyaltyOrderService:      loyaltyOrderService,
		rolesGw:                  rolesGw,
	}, nil
}

func (s *orderService) Create(ctx context.Context, dto *order_dto.CreateOrderDTO) (*order_models.OrderResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	salesSettings, apiErr := s.settingsSalesService.GetSalesSettings(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	order, apiErr := order_utils.ConvertCreateOrderDTOToModel(dto, cfg)
	if apiErr != nil {
		return nil, apiErr
	}

	// Find the pending status from system status types
	for _, status := range salesSettings.OrderStatuses {
		if status.IsSystemStatus && status.SystemStatusType != nil && *status.SystemStatusType == settings_sales_enums.SYSTEM_STATUS_TYPE_PENDING {
			order.Status = status.Title
			break
		}
	}

	// Fallback to first status if no pending status found
	if order.Status == "" && len(salesSettings.OrderStatuses) > 0 {
		order.Status = salesSettings.OrderStatuses[0].Title
	}

	return s.repo.Create(ctx, order)
}

func (s *orderService) GetById(ctx context.Context, id primitive.ObjectID, query string) (*order_models.OrderResponse, *helpers.ApiError) {
	filters, err := order_filters.NewOrderFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}
	return s.repo.GetById(ctx, id, filters)
}

func (s *orderService) GetAll(ctx context.Context, query string, limit int64) ([]order_models.OrderResponse, *helpers.ApiError) {
	filters, filterErr := order_filters.NewOrderFilters(query)
	if filterErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	// Check user branch access
	branchAccess, err := order_utils.CheckUserBranchAccess(ctx, util.GetInjector(ctx))
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to check user branch access",
		}
	}

	// If user has limited access, set the branch filter
	if branchAccess.LimitAccess && branchAccess.BranchId != nil {
		filters.BranchId = branchAccess.BranchId
	}

	pipeline, _ := filters.BuildPipeline(bson.M{}, ctx, false)

	if limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": limit})
	}

	return s.repo.GetAll(ctx, pipeline)
}

func (s *orderService) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*types.DataWithPagination[order_models.OrderResponse], *helpers.ApiError) {
	filters, err := order_filters.NewOrderFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	cfg, _ := util.GetReqAppCfg(ctx)

	userType := cfg.User.Type
	if userType == enums.Vendor {
		filters.VendorId = &cfg.User.Id
	} else if userType == enums.Buyer {
		filters.BuyerId = &cfg.User.Id
	}

	// Check user branch access
	branchAccess, err := order_utils.CheckUserBranchAccess(ctx, util.GetInjector(ctx))
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to check user branch access, error: " + err.Error(),
		}
	}

	// If user has limited access, set the branch filter
	if branchAccess.LimitAccess && branchAccess.BranchId != nil {
		filters.BranchId = branchAccess.BranchId
	}

	return s.repo.GetWithPagination(ctx, skip, limit, *filters)
}

func (s *orderService) UpdateAdminSeen(ctx context.Context, id primitive.ObjectID, dto *order_dto.UpdateAdminSeenDTO) (*order_models.OrderResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	updates := map[string]interface{}{
		"adminSeen": dto.AdminSeen,
	}

	return s.repo.UpdateFieldsById(ctx, id, updates, cfg.User.Id)
}

func (s *orderService) DeleteAllOrders(ctx context.Context) *helpers.ApiError {
	return s.repo.DeleteAllOrders(ctx)
}

func (s *orderService) UpdateTransitionStatus(ctx context.Context, id primitive.ObjectID, status string) (*order_models.OrderResponse, *helpers.ApiError) {
	_, err := s.repo.GetById(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	enumStatus := settings_sales_enums.TransitionStatusType(status)
	if enumStatus == "" {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid status",
		}
	}

	cfg, _ := util.GetReqAppCfg(ctx)

	updates := map[string]interface{}{
		"deliveryInfo.b2bInfo.transitionStatus": enumStatus,
	}
	return s.repo.UpdateFieldsById(ctx, id, updates, cfg.User.Id)

}
