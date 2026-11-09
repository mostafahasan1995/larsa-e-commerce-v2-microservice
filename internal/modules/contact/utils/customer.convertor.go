package contact_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models"
)

func ConvertToCustomerModel(dto *dtos.CreateCustomerDto, additionalOptions types.AdditionalOptions) *models.Customer {
	return &models.Customer{
		BaseCustomer: models.BaseCustomer{
			BasicInformation: models.CustomerBasicInformation{
				FullName:             dto.BasicInformation.FullName,
				OwnerName:            dto.BasicInformation.OwnerName,
				CompanyName:          dto.BasicInformation.CompanyName,
				Service:              dto.BasicInformation.Service,
				CustomsLicenseNumber: dto.BasicInformation.CustomsLicenseNumber,
				BusinessType:         dto.BasicInformation.BusinessType,
				Domain:               dto.BasicInformation.Domain,
			},
			BusinessContactInfo: models.CustomerBusinessContact{
				CommercialRegistrationNumber: dto.BusinessContactInfo.CommercialRegistrationNumber,
				TaxIdentificationNumber:      dto.BusinessContactInfo.TaxIdentificationNumber,
				MobileCountryCode:            dto.BusinessContactInfo.MobileCountryCode,
				MobileNumber:                 dto.BusinessContactInfo.MobileNumber,
				HeadOfficeAddress:            dto.BusinessContactInfo.HeadOfficeAddress,
				AuthorizedCountries:          dto.BusinessContactInfo.AuthorizedCountries,
				WebsiteLink:                  dto.BusinessContactInfo.WebsiteLink,
				Email:                        dto.BusinessContactInfo.Email,
			},
			BankDetails: models.CustomerBankDetails{
				BankName:          dto.BankDetails.BankName,
				AccountNumber:     dto.BankDetails.AccountNumber,
				AccountHolderName: dto.BankDetails.AccountHolderName,
				IBAN:              dto.BankDetails.IBAN,
				SwiftCode:         dto.BankDetails.SwiftCode,
			},
			AccountInformation: models.CustomerAccountInformation{
				Deactivated: dto.AccountInformation.Deactivated,
				Expired:     dto.AccountInformation.Expired,
				Suspended:   dto.AccountInformation.Suspended,
				Live:        dto.AccountInformation.Live,
			},
			NotificationMethods: models.CustomerNotificationMethods{
				Email:     dto.NotificationMethods.Email,
				SMS:       dto.NotificationMethods.SMS,
				Whatsapp:  dto.NotificationMethods.Whatsapp,
				Dashboard: dto.NotificationMethods.Dashboard,
			},
			Subscription: models.CustomerSubscription{
				Plan:               dto.Subscription.Plan,
				StartsAt:           dto.Subscription.StartsAt,
				EndsAt:             dto.Subscription.EndsAt,
				TrialEndsAt:        dto.Subscription.TrialEndsAt,
				UnlimitedDuration:  dto.Subscription.UnlimitedDuration,
				TrialPeriod:        dto.Subscription.TrialPeriod,
				SubscriptionPeriod: dto.Subscription.SubscriptionPeriod,
				SubscriptionAmount: dto.Subscription.SubscriptionAmount,
				DiscountType:       dto.Subscription.DiscountType,
				DiscountAmount:     dto.Subscription.DiscountAmount,
			},
			Attachments:       dto.Attachments,
			InvoiceIds:        dto.InvoiceIds,
			PaymentIds:        dto.PaymentIds,
			TransactionIds:    dto.TransactionIds,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
}
