package quotation_dto

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models/agent"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShippingEmailDto struct {
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

type ShippingAgentDto struct {
	AgentId     *primitive.ObjectID `json:"agentId,omitempty"`
	CountryIso3 string              `bson:"countryIso3,omitempty" json:"countryIso3,omitempty"`
}

type SubmitToAgentDto struct {
	Type        quotation_enums.ShippingType `json:"type" validate:"required,oneof=customs_brokers"`
	From        ShippingAgentDto             `json:"from" validate:"required"`
	To          ShippingAgentDto             `json:"to" validate:"required"`
	Services    []string                     `json:"services"`
	Email       ShippingEmailDto             `json:"email" validate:"required"`
	QuotationId primitive.ObjectID           `json:"quotationId" validate:"required"`
}

type CostDto struct {
	Amount       float64 `json:"amount" validate:"required,min=0"`
	CurrencyCode string  `json:"currencyCode" validate:"required"`
}

type PlatformCostDto struct {
	Amount       float64                          `json:"amount" validate:"required,min=0"`
	CurrencyCode string                           `json:"currencyCode" validate:"required"`
	Type         quotation_enums.PlatformCostType `json:"type" validate:"required,oneof=fixed percentage"`
}

func (dto *SubmitToAgentDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	quotationRepo := base_repos.NewBaseRepo[quotation_models.Quotation, quotation_models.QuotationResponse](args.DB, enums.QUOTATION)
	agentRepo := base_repos.NewBaseRepo[models.Agent, models.AgentResponse](args.DB, enums.AGENT)
	countryRepo := base_repos.NewBaseRepo[config_location_models.Country, config_location_models.CountryResponse](args.DB, enums.COUNTRY)

	// Validate quotation exists and is in correct status
	quotation, apiErr := quotationRepo.GetById(args.Context, dto.QuotationId, []bson.M{})
	if apiErr != nil {
		return nil, apiErr
	}

	uniqeCountries := make(map[string]bool)
	uniqeCountries[dto.From.CountryIso3] = true
	uniqeCountries[dto.To.CountryIso3] = true

	countries, apiErr := countryRepo.GetAll(args.Context, []bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{
						"code.iso3": dto.From.CountryIso3,
					},
					{
						"code.iso3": dto.To.CountryIso3,
					},
				},
			},
		},
	})

	if apiErr != nil {
		return nil, apiErr
	}

	if len(uniqeCountries) != len(countries) {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Some Countries not found",
		}
	}

	if quotation.Status != quotation_enums.PENDING_ADMIN {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Quotation is not pending beebuy",
		}
	}

	// Validate From agent
	if dto.From.AgentId != nil {
		_, apiErr = agentRepo.GetById(args.Context, *dto.From.AgentId, []bson.M{})
		if apiErr != nil {
			return nil, apiErr
		}
	} else {
		agent, _ := agentRepo.FindOne(args.Context, []bson.M{
			{
				"$match": bson.M{
					"businessInfo.authorizedCountry": dto.From.CountryIso3,
				},
			},
		})
		if agent == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "No agent found for the specified country",
			}
		}

		dto.From.AgentId = &agent.Id
	}

	// Validate To agent
	if dto.To.AgentId != nil {
		_, apiErr = agentRepo.GetById(args.Context, *dto.To.AgentId, []bson.M{})
		if apiErr != nil {
			return nil, apiErr
		}
	} else {
		agent, _ := agentRepo.FindOne(args.Context, []bson.M{
			{
				"$match": bson.M{
					"businessInfo.authorizedCountry": dto.To.CountryIso3,
				},
			},
		})
		if agent == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "No agent found for the specified country",
			}
		}

		dto.To.AgentId = &agent.Id
	}

	return nil, nil
}

type UpsertShippingAgentDto struct {
	Method           *quotation_enums.ShippingMethod `bson:"method,omitempty" json:"method,omitempty"`
	CountryIso3      string                          `bson:"countryIso3,omitempty" json:"countryIso3,omitempty"`
	CustomsClearance *CostDto                        `bson:"customsClearance,omitempty" json:"customsClearance,omitempty"`
	Logistic         *CostDto                        `bson:"logistic,omitempty" json:"logistic,omitempty"`
	Insurance        *CostDto                        `bson:"insurance,omitempty" json:"insurance,omitempty"`
}

func (dto *UpsertShippingAgentDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	shippingId := chi.URLParam(args.Request, "id")
	shippingObjectId, _ := primitive.ObjectIDFromHex(shippingId)

	shippingRepo := base_repos.NewBaseRepo[quotation_models.Shipping, quotation_models.ShippingResponse](args.DB, enums.SHIPPING)

	shipping, apiErr := shippingRepo.GetById(args.Context, shippingObjectId, []bson.M{})
	if apiErr != nil {
		return nil, apiErr
	}

	if shipping.Status != quotation_enums.SHIPPING_PENDING {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Shipping is not pending",
		}
	}

	return []middleware.CustomValidatorContext{
		{
			Key:   string(quotation_enums.SHIPPING_CONTEXT_KEY),
			Value: shipping,
		},
	}, nil
}

type SubmitToVendorDto struct {
	PlatformCost PlatformCostDto `json:"platformCost" validate:"required"`
}

func (dto *SubmitToVendorDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	shippingId := chi.URLParam(args.Request, "id")
	shippingObjectId, _ := primitive.ObjectIDFromHex(shippingId)

	shippingRepo := base_repos.NewBaseRepo[quotation_models.Shipping, quotation_models.ShippingResponse](args.DB, enums.SHIPPING)

	shipping, apiErr := shippingRepo.GetById(args.Context, shippingObjectId, []bson.M{})
	if apiErr != nil {
		return nil, apiErr
	}

	if shipping.Status != quotation_enums.SHIPPING_APPROVED {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Shipping is not approved",
		}
	}

	quotationRepo := base_repos.NewBaseRepo[quotation_models.Quotation, quotation_models.QuotationResponse](args.DB, enums.QUOTATION)
	quotation, apiErr := quotationRepo.GetById(args.Context, shipping.QuotationId, []bson.M{})
	if apiErr != nil {
		return nil, apiErr
	}

	if quotation.Status != quotation_enums.PENDING_ADMIN {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Quotation is not pending beebuy",
		}
	}

	return []middleware.CustomValidatorContext{
		{
			Key:   string(quotation_enums.SHIPPING_CONTEXT_KEY),
			Value: shipping,
		},
	}, nil
}
