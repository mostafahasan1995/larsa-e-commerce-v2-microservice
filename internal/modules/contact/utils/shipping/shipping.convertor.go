package contact_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos/shipping"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models/shipping"
)

func ConvertToShippingModel(dto *dtos.CreateShippingDto, additionalOptions types.AdditionalOptions) *models.Shipping {
	return &models.Shipping{
		BaseShipping: models.BaseShipping{
			BasicInformation: models.BasicInformation{
				Name:                         dto.BasicInformation.Name,
				Type:                         dto.BasicInformation.Type,
				CommercialRegistrationNumber: dto.BasicInformation.CommercialRegistrationNumber,
				TaxIdentificationNumber:      dto.BasicInformation.TaxIdentificationNumber,
			},
			Service: models.ShippingService{
				LandFreight:                dto.Service.LandFreight,
				AirFreight:                 dto.Service.AirFreight,
				SeaFreight:                 dto.Service.SeaFreight,
				WarehousingAndDistribution: dto.Service.WarehousingAndDistribution,
			},
			Contract: models.ShippingContract{
				Email:             dto.Contract.Email,
				HeadOfficeAddress: dto.Contract.HeadOfficeAddress,
				Phone:             dto.Contract.Phone,
				CountryCodes:      dto.Contract.CountryCodes,
				Website:           dto.Contract.Website,
			},
			Financial: models.ShippingFinancial{
				StartDate:       dto.Financial.StartDate,
				EndDate:         dto.Financial.EndDate,
				BankInformation: dto.Financial.BankInformation,
				Attachments:     dto.Financial.Attachments,
			},
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
}
