package quotation_service

import (
	"context"
	"fmt"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	profile_services "larsa-e-commerce-microservice-v2/internal/modules/profile/services"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_filters "larsa-e-commerce-microservice-v2/internal/modules/quotation/filters"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	quotation_repos "larsa-e-commerce-microservice-v2/internal/modules/quotation/repo"
	quotation_utils "larsa-e-commerce-microservice-v2/internal/modules/quotation/utils"

	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuotationAdminService interface {
	UpdateQuotationStatus(ctx context.Context, quotationID string, status quotation_enums.QuotationStatus) *helpers.ApiError
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError)
	GetById(ctx context.Context, id string, query string) (*quotation_models.QuotationResponse, *helpers.ApiError)
	SetAdminFee(ctx context.Context, quotationID string, dto quotation_dto.CostDto) (*quotation_models.QuotationResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, quotationID string) *helpers.ApiError
	DeleteAll(ctx context.Context) *helpers.ApiError
	// Agent & Shipping Company
	SetShippingCost(ctx context.Context, quotationID string, dto quotation_dto.SetShippingCostDTO) (*quotation_models.QuotationResponse, *helpers.ApiError)
	SetAgentFee(ctx context.Context, quotationID string, dto quotation_dto.SetAgentFeeDTO) (*quotation_models.QuotationResponse, *helpers.ApiError)
	ConvertQuotationToOrder(ctx context.Context, quotationId primitive.ObjectID) (*order_models.Order, *helpers.ApiError)
}

type quotationAdmin struct {
	repo                           quotation_repos.QuotationRepo
	profileService                 profile_services.ProfileService
	quotationOrderConvertorService QuotationOrderConvertorService
}

func NewQuotationAdminService(i *do.Injector) (QuotationAdminService, error) {
	return &quotationAdmin{
		repo:                           do.MustInvoke[quotation_repos.QuotationRepo](i),
		profileService:                 do.MustInvoke[profile_services.ProfileService](i),
		quotationOrderConvertorService: do.MustInvoke[QuotationOrderConvertorService](i),
	}, nil
}

func (s *quotationAdmin) GetById(ctx context.Context, id string, query string) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	objectId, parseErr := primitive.ObjectIDFromHex(id)
	if parseErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid ID format",
		}
	}
	filter, err := quotation_filters.NewQuotationFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}
	return s.repo.GetById(ctx, objectId, *filter)
}

func (s *quotationAdmin) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError) {
	config, _ := util.GetReqAppCfg(ctx)
	filter, err := quotation_filters.NewQuotationFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	if config != nil && config.User != nil {
		if config.User.Type == enums.Agent || config.User.Type == enums.ShippingCompany {
			filter.AnyPartyUserId = &config.User.Id
		} else if config.User.Type == enums.Vendor {
			filter.VendorProfileId = &config.User.Id
		}
	}

	return s.repo.GetWithPagination(ctx, skip, limit, *filter)
}

func (s *quotationAdmin) UpdateQuotationStatus(ctx context.Context, quotationId string, status quotation_enums.QuotationStatus) *helpers.ApiError {
	quotation, apiErr := s.GetById(ctx, quotationId, `{"popFields":["quotationBuyerInfo.buyerProfileId","vendorProfileId"]}`)
	if apiErr != nil {
		return apiErr
	}

	//TODO: We removed the validation for now, just for testing, please uncomment this when you are done testing.
	// config, _ := util.GetReqAppCfg(ctx)

	// if !quotation_utils.IsUpdateQuotationStatusValid(quotation, status) {
	// 	return &helpers.ApiError{
	// 		StatusCode: 400,
	// 		Message:    "Invalid status",
	// 	}
	// }

	// if !quotation_utils.IsUserAllowedToUpdateQuotationStatus(quotation, status, config.User) {
	// 	return &helpers.ApiError{
	// 		StatusCode: 400,
	// 		Message:    "User is not allowed to update quotation status",
	// 	}
	// }

	historyEntry := quotation_utils.GetHistoryEntry(quotation, status)
	newHistory := append(quotation.History, *historyEntry)

	s.repo.UpdateQuotationStatus(ctx, quotation.Id, status, newHistory)

	if status == quotation_enums.CONVERT_TO_ORDER {
		order, convertErr := s.ConvertQuotationToOrder(ctx, quotation.Id)
		if convertErr != nil {
			return convertErr
		}

		orderRs, orderErr := s.repo.CreateOrderFromQuotation(ctx, quotation.Id, *order)
		if orderErr != nil {
			return orderErr
		}

		fmt.Println("orderRs", orderRs)
	} else {
		//TODO: handler other status changes to send emails and notifications.
	}
	return nil
}

func (s *quotationAdmin) SetAdminFee(ctx context.Context, quotationID string, dto quotation_dto.CostDto) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	quotation, apiErr := s.GetById(ctx, quotationID, "")
	if apiErr != nil {
		return nil, apiErr
	}
	update := map[string]interface{}{
		"fee.superAdminFee": quotation_models.QuotationPrice{Amount: dto.Amount, CurrencyCode: dto.CurrencyCode},
	}
	if len(update) == 0 {
		return quotation, nil
	}

	if quotation.Status != quotation_enums.PENDING_ADMIN {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Can't set admin fee for this quotation, due to the quotation is not in pending admin status"}
	}

	config, _ := util.GetReqAppCfg(ctx)
	_, err := s.repo.UpdateFieldsById(ctx, quotation.Id, update, config.User.Id)
	if err != nil {
		return nil, err
	}

	// Update history
	additionalInfo := "Admin fee updated by super admin"
	historyEntry := quotation_utils.CreateHistoryEntry(quotation_enums.QUOTATION_UPDATED, additionalInfo)
	newHistory := append(quotation.History, historyEntry)

	_, historyErr := s.repo.UpdateFieldsById(ctx, quotation.Id, map[string]interface{}{"history": newHistory}, config.User.Id)
	if historyErr != nil {
		fmt.Println("Error updating history:", historyErr)
	}

	newQuotation, err := s.GetById(ctx, quotation.Id.Hex(), "")
	if err != nil {
		return nil, err
	}
	return newQuotation, nil
}

func (s *quotationAdmin) DeleteById(ctx context.Context, quotationID string) *helpers.ApiError {
	objectId, parseErr := primitive.ObjectIDFromHex(quotationID)
	if parseErr != nil {
		return &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid ID format",
		}
	}

	return s.repo.DeleteById(ctx, objectId, nil)
}

func (s *quotationAdmin) ConvertQuotationToOrder(ctx context.Context, quotationId primitive.ObjectID) (*order_models.Order, *helpers.ApiError) {
	quotation, apiErr := s.GetById(ctx, quotationId.Hex(), "")
	if apiErr != nil {
		return nil, apiErr
	}
	return s.quotationOrderConvertorService.ConvertQuotationToOrder(ctx, quotation)
}

func (s *quotationAdmin) DeleteAll(ctx context.Context) *helpers.ApiError {
	return s.repo.DeleteAll(ctx)
}
