package contact_dtos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	payment_models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerBasicInformationDto struct {
	FullName             types.Localizable[string] `json:"fullName" validate:"required,localizable"`
	OwnerName            types.Localizable[string] `json:"ownerName" validate:"required,localizable"`
	CompanyName          types.Localizable[string] `json:"companyName" validate:"required,localizable"`
	Service              string                    `json:"service" validate:"required"`
	CustomsLicenseNumber string                    `json:"customsLicenseNumber" validate:"required"`
	BusinessType         string                    `json:"businessType" validate:"required"`
	Domain               string                    `json:"domain" validate:"required"`
}

type CustomerAccountInformationDto struct {
	Deactivated bool `json:"deactivated"`
	Expired     bool `json:"expired"`
	Suspended   bool `json:"suspended"`
	Live        bool `json:"live"`
}

type CustomerNotificationMethodsDto struct {
	Email     bool `json:"email"`
	SMS       bool `json:"sms"`
	Whatsapp  bool `json:"whatsapp"`
	Dashboard bool `json:"dashboard"`
}

type CustomerSubscriptionDto struct {
	//TODO with plan id
	Plan              string     `json:"plan" validate:"required"`
	StartsAt          time.Time  `json:"startsAt" validate:"required"`
	EndsAt            *time.Time `json:"endsAt"`
	TrialEndsAt       time.Time  `json:"trialEndsAt" validate:"required"`
	UnlimitedDuration bool       `json:"unlimitedDuration"`
	//by day
	TrialPeriod        int     `json:"trialPeriod" validate:"required"`
	SubscriptionPeriod int     `json:"subscriptionPeriod" validate:"required"`
	SubscriptionAmount int     `json:"subscriptionAmount" validate:"required"`
	DiscountType       string  `json:"discountType" validate:"omitempty,oneof=fixed percentage"`
	DiscountAmount     float64 `json:"discountAmount"`
}

type CustomerBusinessContactDto struct {
	CommercialRegistrationNumber string   `json:"commercialRegistrationNumber" validate:"required"`
	TaxIdentificationNumber      string   `json:"taxIdentificationNumber" validate:"required"`
	MobileCountryCode            string   `json:"mobileCountryCode" validate:"required"`
	MobileNumber                 string   `json:"mobileNumber" validate:"required"`
	HeadOfficeAddress            string   `json:"headOfficeAddress" validate:"required"`
	AuthorizedCountries          []string `json:"authorizedCountries" validate:"required"`
	WebsiteLink                  string   `json:"websiteLink" validate:"required,url"`
	Email                        string   `json:"email" validate:"required,email"`
	Password                     *string  `json:"password" validate:"omitempty,min=8"`
}

type CustomerBankDetailsDto struct {
	BankName          string `json:"bankName" validate:"required"`
	AccountNumber     string `json:"accountNumber" validate:"required"`
	AccountHolderName string `json:"accountHolderName" validate:"required"`
	IBAN              string `json:"iban" validate:"required"`
	SwiftCode         string `json:"swiftCode" validate:"required"`
}

type CreateCustomerDto struct {
	UserId              *primitive.ObjectID
	BasicInformation    CustomerBasicInformationDto    `json:"basicInformation" validate:"required"`
	BusinessContactInfo CustomerBusinessContactDto     `json:"businessContactInfo" validate:"required"`
	Subscription        CustomerSubscriptionDto        `json:"subscription" validate:"required"`
	BankDetails         CustomerBankDetailsDto         `json:"bankDetails" validate:"required"`
	NotificationMethods CustomerNotificationMethodsDto `json:"notificationMethods" validate:"omitempty"`
	AccountInformation  CustomerAccountInformationDto  `json:"accountInformation" validate:"omitempty"`
	Attachments         []types.FileFieldDto           `json:"attachments" validate:"omitempty,dive"`

	InvoiceIds     []primitive.ObjectID `json:"invoiceIds" validate:"omitempty"`
	PaymentIds     []primitive.ObjectID `json:"paymentIds" validate:"omitempty"`
	TransactionIds []primitive.ObjectID `json:"transactionIds" validate:"omitempty"`
}

// validateIds is a helper function to validate IDs using a repository
func validateIds[M any, R any](ctx context.Context, ids []primitive.ObjectID, repo *base_repos.BaseRepo[M, R], errChan chan<- *helpers.ApiError, cancel context.CancelFunc) {
	select {
	case <-ctx.Done():
		return
	default:
		if len(ids) > 0 {
			_, apiErr := repo.GetByIds(ctx, ids, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
			if apiErr != nil {
				errChan <- apiErr
				cancel() // Cancel other goroutines if this one fails
			}
		} else {
			errChan <- nil
		}
	}
}

func (dto *CreateCustomerDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	invoiceRepo := base_repos.NewBaseRepo[payment_models.Invoice, payment_models.InvoiceResponse](args.DB, enums.INVOICE)
	paymentRepo := base_repos.NewBaseRepo[payment_models.Payment, payment_models.PaymentResponse](args.DB, enums.PAYMENT)
	transactionRepo := base_repos.NewBaseRepo[payment_models.Transaction, payment_models.TransactionResponse](args.DB, enums.TRANSACTION)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(args.Context)
	defer cancel()

	errChan := make(chan *helpers.ApiError, 3)
	done := make(chan bool)

	// Launch goroutines for each validation
	go validateIds[payment_models.Invoice, payment_models.InvoiceResponse](ctx, dto.InvoiceIds, invoiceRepo, errChan, cancel)
	go validateIds[payment_models.Payment, payment_models.PaymentResponse](ctx, dto.PaymentIds, paymentRepo, errChan, cancel)
	go validateIds[payment_models.Transaction, payment_models.TransactionResponse](ctx, dto.TransactionIds, transactionRepo, errChan, cancel)

	// Wait for either an error or all goroutines to complete
	select {
	case err := <-errChan:
		return nil, err
	case <-done:
		return nil, nil
	}
}
