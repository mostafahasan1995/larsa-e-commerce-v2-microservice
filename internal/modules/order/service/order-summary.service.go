package order_service

import (
	"context"
	"fmt"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"

	excel_services "larsa-e-commerce-microservice-v2/internal/common/services/excel"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_filters "larsa-e-commerce-microservice-v2/internal/modules/order/filters"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_repos "larsa-e-commerce-microservice-v2/internal/modules/order/repos"
	order_utils "larsa-e-commerce-microservice-v2/internal/modules/order/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderSummaryService interface {
	GetUpcomingShipments(ctx context.Context, limit int64) (*types.DataWithPagination[order_models.OrderResponse], *helpers.ApiError)
	GetRecentActivity(ctx context.Context, branchId *primitive.ObjectID, query string) (map[string]float64, *helpers.ApiError)
	GetBuyerActivity(ctx context.Context, branchId *primitive.ObjectID, query string) (*order_dto.BuyerActivityResponse, *helpers.ApiError)
	ExportOrders(ctx context.Context, query string, limit int64) ([]byte, *helpers.ApiError)
}

type orderSummaryService struct {
	repo         order_repos.OrderRepo
	excelService excel_services.ExcelService
	orderService OrderService
}

func NewOrderSummaryService(i *do.Injector) (OrderSummaryService, error) {
	repo := do.MustInvoke[order_repos.OrderRepo](i)
	excelService := do.MustInvoke[excel_services.ExcelService](i)
	orderService := do.MustInvoke[OrderService](i)
	return &orderSummaryService{
		repo:         repo,
		excelService: excelService,
		orderService: orderService,
	}, nil
}

func (s *orderSummaryService) GetUpcomingShipments(ctx context.Context, limit int64) (*types.DataWithPagination[order_models.OrderResponse], *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	filters := order_filters.OrderFilters{}
	userType := cfg.User.Type
	if userType == enums.Vendor {
		filters.VendorId = &cfg.User.Id
	} else if userType == enums.Buyer {
		filters.BuyerId = &cfg.User.Id
	}
	filters.Sorts = []base_filters.Sort{
		{
			Field: "createdAt",
			Order: base_filters.Desc,
		},
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

	return s.repo.GetWithPagination(ctx, 0, limit, filters)
}

func (s *orderSummaryService) GetRecentActivity(ctx context.Context, branchId *primitive.ObjectID, query string) (map[string]float64, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	filters, err := order_filters.NewOrderFilters(query)
	if err != nil {
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
			Message:    "Failed to check user branch access, error: " + err.Error(),
		}
	}

	// If user has limited access, set the branch filter
	if branchAccess.LimitAccess && branchAccess.BranchId != nil {
		filters.BranchId = branchAccess.BranchId
	}

	return s.repo.GetRecentActivity(ctx, cfg.User.Id, cfg.User.Type, branchId, *filters)
}

func (s *orderSummaryService) GetBuyerActivity(ctx context.Context, branchId *primitive.ObjectID, query string) (*order_dto.BuyerActivityResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	filters, err := order_filters.NewOrderFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	buyerId := cfg.User.Id

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

	return s.repo.GetBuyerActivity(ctx, buyerId, branchId, *filters)
}

func (s *orderSummaryService) ExportOrders(ctx context.Context, query string, limit int64) ([]byte, *helpers.ApiError) {
	// Get all orders
	orders, apiErr := s.orderService.GetAll(ctx, query, limit)
	if apiErr != nil {
		return nil, apiErr
	}

	// Convert orders to export data
	exportData := make([]map[string]interface{}, len(orders))
	for i, order := range orders {
		exportData[i] = order_utils.ConvertOrderToExportData(order)
	}

	// Define Excel columns for order export
	columns := order_utils.GetOrderExportColumns()

	// Generate filename with current date
	now := time.Now()
	filename := fmt.Sprintf("Orders_%s", now.Format("2006-01-02"))

	// Create Excel service and export
	excelData, apiErr := s.excelService.Export(ctx, exportData, columns, filename)
	if apiErr != nil {
		return nil, apiErr
	}

	return excelData, nil
}
