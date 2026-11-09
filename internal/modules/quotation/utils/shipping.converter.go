package quotation_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
)

func ConvertToSubmitToAgentModel(dto *quotation_dto.SubmitToAgentDto, additionalOptions types.AdditionalOptions) *quotation_models.Shipping {

	shipping := &quotation_models.Shipping{
		BaseShipping: quotation_models.BaseShipping{
			Status: quotation_enums.SHIPPING_PENDING,
			From: quotation_models.ShippingAgent{
				AgentId:     *dto.From.AgentId,
				CountryIso3: dto.From.CountryIso3,
			},
			To: quotation_models.ShippingAgent{
				AgentId:     *dto.To.AgentId,
				CountryIso3: dto.To.CountryIso3,
			},
			Type:     dto.Type,
			Services: dto.Services,
			Email: &quotation_models.ShippingEmail{
				Subject: dto.Email.Subject,
				Body:    dto.Email.Body,
			},
			QuotationId: dto.QuotationId,

			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
	return shipping
}

func ConvertToUpsertShippingAgentModel(dto *quotation_dto.UpsertShippingAgentDto, additionalOptions types.AdditionalOptions) *quotation_models.ShippingAgent {
	shippingAgent := &quotation_models.ShippingAgent{
		Method:      dto.Method,
		CountryIso3: dto.CountryIso3,
	}

	if dto.CustomsClearance != nil {
		shippingAgent.CustomsClearance = &quotation_models.Cost{
			Amount:       dto.CustomsClearance.Amount,
			CurrencyCode: dto.CustomsClearance.CurrencyCode,
		}
	}

	if dto.Logistic != nil {
		shippingAgent.Logistic = &quotation_models.Cost{
			Amount:       dto.Logistic.Amount,
			CurrencyCode: dto.Logistic.CurrencyCode,
		}
	}

	if dto.Insurance != nil {
		shippingAgent.Insurance = &quotation_models.Cost{
			Amount:       dto.Insurance.Amount,
			CurrencyCode: dto.Insurance.CurrencyCode,
		}
	}

	return shippingAgent
}

func ConvertToSubmitToVendorModel(dto *quotation_dto.SubmitToVendorDto, additionalOptions types.AdditionalOptions) *quotation_models.PlatformCost {

	platformCost := &quotation_models.PlatformCost{
		Amount:       dto.PlatformCost.Amount,
		CurrencyCode: dto.PlatformCost.CurrencyCode,
		Type:         dto.PlatformCost.Type,
	}

	return platformCost
}
