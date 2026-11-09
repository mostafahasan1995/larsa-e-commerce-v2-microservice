package contact_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models"
)

func ConvertToBrokerModel(dto *dtos.CreateBrokerDto, additionalOptions types.AdditionalOptions) *models.Broker {
	return &models.Broker{
		BaseBroker: models.BaseBroker{
			BasicInformation: models.BrokerBasicInformation{
				FullName:                     dto.BasicInformation.FullName,
				CompanyName:                  dto.BasicInformation.CompanyName,
				CustomsLicenseNumber:         dto.BasicInformation.CustomsLicenseNumber,
				CommercialRegistrationNumber: dto.BasicInformation.CommercialRegistrationNumber,
				TaxIdentificationNumber:      dto.BasicInformation.TaxIdentificationNumber,
				AuthorizedCountry:            dto.BasicInformation.AuthorizedCountry,
				YearsOfExperience:            dto.BasicInformation.YearsOfExperience,
			},
			BankAccount: models.BankAccount{
				BankName:          dto.BankAccount.BankName,
				AccountNumber:     dto.BankAccount.AccountNumber,
				AccountHolderName: dto.BankAccount.AccountHolderName,
				IBAN:              dto.BankAccount.IBAN,
				SWIFTCode:         dto.BankAccount.SWIFTCode,
			},
			Financial: models.BrokerFinancial{
				StartDate:   dto.Financial.StartDate,
				EndDate:     dto.Financial.EndDate,
				Attachments: dto.Financial.Attachments,
			},
			Contract: models.BrokerContract{
				HeadOfficeAddress: dto.Contract.HeadOfficeAddress,
				Email:             dto.Contract.Email,
				PhoneNumber:       dto.Contract.PhoneNumber,
			},
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
}
