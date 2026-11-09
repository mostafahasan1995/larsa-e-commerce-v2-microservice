package contact_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos/agent"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models/agent"
)

func ConvertToAgentModel(dto *dtos.CreateAgentDto, additionalOptions types.AdditionalOptions) *models.Agent {
	return &models.Agent{
		BaseAgent: models.BaseAgent{
			ResponsibleForLogistics: dto.ResponsibleForLogistics,
			Active:                  dto.Active,
			BasicInformation: models.AgentBasicInfo{
				FirstName:         dto.BasicInformation.FirstName,
				LastName:          dto.BasicInformation.LastName,
				CompanyName:       dto.BasicInformation.CompanyName,
				LicenseNumber:     dto.BasicInformation.LicenseNumber,
				YearsOfExperience: dto.BasicInformation.YearsOfExperience,
			},
			BusinessInfo: models.AgentBusinessInfo{
				CommercialRegistrationNumber: dto.BusinessInfo.CommercialRegistrationNumber,
				TaxIdentificationNumber:      dto.BusinessInfo.TaxIdentificationNumber,
				PhoneNumber:                  dto.BusinessInfo.PhoneNumber,
				HeadOfficeAddress:            dto.BusinessInfo.HeadOfficeAddress,
				OfficialEmail:                dto.BusinessInfo.OfficialEmail,
				AuthorizedCountryCodes:       dto.BusinessInfo.AuthorizedCountryCodes,
				WebsiteLink:                  dto.BusinessInfo.WebsiteLink,
			},
			FinancialInfo: models.AgentFinancialInfo{
				ShippingProfitPercentage:     dto.FinancialInfo.ShippingProfitPercentage,
				SubscriptionProfitPercentage: dto.FinancialInfo.SubscriptionProfitPercentage,
			},
			ContractInfo: models.AgentContractInfo{
				StartDate:   dto.ContractInfo.StartDate,
				EndDate:     dto.ContractInfo.EndDate,
				Attachments: dto.ContractInfo.Attachments,
			},
			BankAccountInfo: models.AgentBankAccountInfo{
				BankName:          dto.BankAccountInfo.BankName,
				AccountNumber:     dto.BankAccountInfo.AccountNumber,
				AccountHolderName: dto.BankAccountInfo.AccountHolderName,
				IBAN:              dto.BankAccountInfo.Iban,
				SWIFTCode:         dto.BankAccountInfo.SwiftCode,
			},
			Attachments:       dto.Attachments,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
}
