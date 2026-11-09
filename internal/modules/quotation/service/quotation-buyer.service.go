package quotation_service

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_service "larsa-e-commerce-microservice-v2/internal/modules/order/service"
	profile_services "larsa-e-commerce-microservice-v2/internal/modules/profile/services"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_filters "larsa-e-commerce-microservice-v2/internal/modules/quotation/filters"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	quotation_repos "larsa-e-commerce-microservice-v2/internal/modules/quotation/repo"
	quotation_utils "larsa-e-commerce-microservice-v2/internal/modules/quotation/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuotationBuyerService interface {
	CreateRequestQuotation(ctx context.Context, dto *quotation_dto.CreateQuotationDTO) (*quotation_models.QuotationResponse, *helpers.ApiError)
	GetResponseQuotations(ctx context.Context, skip, limit int64, requestQuotationId string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError)
	GetById(ctx context.Context, id string, query string) (*quotation_models.QuotationResponse, *helpers.ApiError)
	GetGeneralQuotations(ctx context.Context, skip, limit int64, query string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError)
	RequestGeneralQuotation(ctx context.Context, id string, buyerProfileDto quotation_dto.QuotationBuyerInfoDTO) (*quotation_models.QuotationResponse, *helpers.ApiError)
}

type quotationBuyer struct {
	repo                           quotation_repos.QuotationRepo
	orderService                   order_service.OrderService
	profileService                 profile_services.ProfileService
	quotationOrderConvertorService QuotationOrderConvertorService
}

func NewQuotationBuyerService(i *do.Injector) (QuotationBuyerService, error) {
	return &quotationBuyer{
		repo:                           do.MustInvoke[quotation_repos.QuotationRepo](i),
		orderService:                   do.MustInvoke[order_service.OrderService](i),
		profileService:                 do.MustInvoke[profile_services.ProfileService](i),
		quotationOrderConvertorService: do.MustInvoke[QuotationOrderConvertorService](i),
	}, nil
}

func (s *quotationBuyer) CreateRequestQuotation(ctx context.Context, dto *quotation_dto.CreateQuotationDTO) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	quotation, apiErr := quotation_utils.ConvertCreateQuotationDTOToModel(dto, ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	config, _ := util.GetReqAppCfg(ctx)
	quotation.Code = quotation.GenerateCode()
	quotation.Type = quotation_enums.REQUEST
	quotation.Status = quotation_enums.PENDING_VENDOR
	quotation.QuotationBuyerInfo.BuyerProfileId = config.User.Id

	// Validate that ProductsObjects has the same length as Products
	if len(dto.ProductsObjects) != len(dto.Products) {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Products validation failed: mismatch between requested products and available products",
		}
	}

	for i := range dto.Products {
		//TODO: need to set a default value for price in the quotation product from the product object.
		productObject := dto.ProductsObjects[i]

		if len(productObject.Prices) == 0 {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    "Product has no prices available",
			}
		}

		if len(productObject.Prices[0].BasePrice.Items) == 0 {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    "Product price has no items available",
			}
		}
		quotation.Products[i].Price.Amount = productObject.Prices[0].BasePrice.Items[0].Value
		quotation.Products[i].Price.CurrencyCode = productObject.Prices[0].BasePrice.Items[0].CurrencyCode

		// Check if discount exists before accessing its value
		if productObject.Prices[0].BasePrice.Discount != nil {
			quotation.Products[i].DiscountPercent = productObject.Prices[0].BasePrice.Discount.Value
		} else {
			quotation.Products[i].DiscountPercent = 0.0 // Default to no discount
		}
	}

	result, apiErr := s.repo.Create(ctx, quotation)
	if apiErr != nil {
		return nil, apiErr
	}
	return result, nil
}

func (s *quotationBuyer) GetById(ctx context.Context, id string, query string) (*quotation_models.QuotationResponse, *helpers.ApiError) {
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

// Get response quotations for a request quotation
func (s *quotationBuyer) GetResponseQuotations(ctx context.Context, skip, limit int64, requestQuotationId string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError) {
	requestQuotationIdObj, err := primitive.ObjectIDFromHex(requestQuotationId)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid request quotation ID",
		}
	}
	filter := quotation_filters.QuotationFilters{
		RequestQuotationId: &requestQuotationIdObj,
	}

	return s.repo.GetWithPagination(ctx, skip, limit, filter)
}

// This function cover: get my requests quotations and get general quotations.
// it do all of that using the quotation filters.
func (s *quotationBuyer) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError) {
	filter, filterError := quotation_filters.NewQuotationFilters(query)
	if filterError != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	config, _ := util.GetReqAppCfg(ctx)
	filter.BuyerProfileId = &config.User.Id

	return s.repo.GetWithPagination(ctx, skip, limit, *filter)
}

func (s *quotationBuyer) GetGeneralQuotations(ctx context.Context, skip, limit int64, query string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError) {
	filter, filterError := quotation_filters.NewQuotationFilters(query)
	if filterError != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	generalType := quotation_enums.GENERAL
	filter.Type = &generalType

	approvedStatus := quotation_enums.APPROVED
	filter.Status = &approvedStatus

	return s.repo.GetWithPagination(ctx, skip, limit, *filter)
}

func (s *quotationBuyer) RequestGeneralQuotation(ctx context.Context, id string, buyerProfileDto quotation_dto.QuotationBuyerInfoDTO) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	quotationId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid quotation ID",
		}
	}

	quotation, findErr := s.repo.GetById(ctx, quotationId, quotation_filters.QuotationFilters{})
	if findErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Quotation not found",
		}
	}

	if quotation.Type != quotation_enums.GENERAL {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Quotation is not a general quotation",
		}
	}

	if quotation.Status != quotation_enums.APPROVED {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Quotation is not approved",
		}
	}

	buyerInfo := quotation_utils.ConvertQuotationBuyerInfoDTOToModel(&buyerProfileDto)
	quotation.QuotationBuyerInfo = quotation_models.QuotationBuyerInfoResponse{
		BaseQuotationBuyerInfo: buyerInfo.BaseQuotationBuyerInfo,
	}

	//Create a new quotation with the same data but with the status requested general quotation
	quotation.GeneralQuotationId = &quotationId

	quotation.Id = primitive.NewObjectID()

	newQuotation := quotation_models.Quotation{
		BaseQuotation:      quotation.BaseQuotation,
		QuotationBuyerInfo: buyerInfo,
	}

	return s.repo.Create(ctx, newQuotation)
}
