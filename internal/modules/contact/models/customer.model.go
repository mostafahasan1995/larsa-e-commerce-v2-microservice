package contact_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	payment_models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerBasicInformation struct {
	FullName             types.Localizable[string] `json:"fullName"`
	OwnerName            types.Localizable[string] `json:"ownerName"`
	CompanyName          types.Localizable[string] `json:"companyName"`
	Service              string                    `json:"service"`
	CustomsLicenseNumber string                    `json:"customsLicenseNumber"`
	BusinessType         string                    `json:"businessType"`
	Domain               string                    `json:"domain"`
}

type CustomerAccountInformation struct {
	Deactivated bool `json:"deactivated"`
	Expired     bool `json:"expired"`
	Suspended   bool `json:"suspended"`
	Live        bool `json:"live"`
}

type CustomerBusinessContact struct {
	CommercialRegistrationNumber string   `bson:"commercialRegistrationNumber" json:"commercialRegistrationNumber"`
	TaxIdentificationNumber      string   `bson:"taxIdentificationNumber" json:"taxIdentificationNumber"`
	MobileCountryCode            string   `bson:"mobileCountryCode" json:"mobileCountryCode"`
	MobileNumber                 string   `bson:"mobileNumber" json:"mobileNumber"`
	HeadOfficeAddress            string   `bson:"headOfficeAddress" json:"headOfficeAddress"`
	AuthorizedCountries          []string `bson:"authorizedCountries" json:"authorizedCountries"`
	WebsiteLink                  string   `bson:"websiteLink" json:"websiteLink"`
	Email                        string   `bson:"email" json:"email"`
}

type CustomerBankDetails struct {
	BankName          string `bson:"bankName" json:"bankName"`
	AccountNumber     string `bson:"accountNumber" json:"accountNumber"`
	AccountHolderName string `bson:"accountHolderName" json:"accountHolderName"`
	IBAN              string `bson:"iban" json:"iban"`
	SwiftCode         string `bson:"swiftCode" json:"swiftCode"`
}

type CustomerNotificationMethods struct {
	Email     bool `bson:"email" json:"email"`
	SMS       bool `bson:"sms" json:"sms"`
	Whatsapp  bool `bson:"whatsapp" json:"whatsapp"`
	Dashboard bool `bson:"dashboard" json:"dashboard"`
}

type CustomerSubscription struct {
	Plan               string     `bson:"plan" json:"plan"`
	StartsAt           time.Time  `bson:"startsAt" json:"startsAt"`
	EndsAt             *time.Time `bson:"endsAt,omitempty" json:"endsAt,omitempty"`
	TrialEndsAt        time.Time  `bson:"trialEndsAt" json:"trialEndsAt"`
	UnlimitedDuration  bool       `bson:"unlimitedDuration" json:"unlimitedDuration"`
	TrialPeriod        int        `bson:"trialPeriod" json:"trialPeriod"`
	SubscriptionPeriod int        `bson:"subscriptionPeriod" json:"subscriptionPeriod"`
	SubscriptionAmount int        `bson:"subscriptionAmount" json:"subscriptionAmount"`
	DiscountType       string     `bson:"discountType" json:"discountType"`
	DiscountAmount     float64    `bson:"discountAmount" json:"discountAmount"`
}

type BaseCustomer struct {
	BasicInformation    CustomerBasicInformation    `bson:"basicInformation" json:"basicInformation"`
	BusinessContactInfo CustomerBusinessContact     `bson:"businessContactInfo" json:"businessContactInfo"`
	BankDetails         CustomerBankDetails         `bson:"bankDetails" json:"bankDetails"`
	AccountInformation  CustomerAccountInformation  `bson:"accountInformation" json:"accountInformation"`
	NotificationMethods CustomerNotificationMethods `bson:"notificationMethods" json:"notificationMethods"`
	Subscription        CustomerSubscription        `bson:"subscription" json:"subscription"`
	Attachments         []types.FileFieldDto        `bson:"attachments" json:"attachments"`
	InvoiceIds          []primitive.ObjectID        `bson:"invoiceIds" json:"invoiceIds"`
	PaymentIds          []primitive.ObjectID        `bson:"paymentIds" json:"paymentIds"`
	TransactionIds      []primitive.ObjectID        `bson:"transactionIds" json:"transactionIds"`

	base_models.MetaDataBaseModel `bson:",inline"`
}
type Customer struct {
	BaseCustomer `bson:",inline"`
}

type CustomerResponse struct {
	BaseCustomer `bson:",inline"`
	Invoices     []payment_models.InvoiceResponse     `bson:"invoices" json:"invoices"`
	Payments     []payment_models.PaymentResponse     `bson:"payments" json:"payments"`
	Transactions []payment_models.TransactionResponse `bson:"transactions" json:"transactions"`
}

type CustomersWithPagination = types.DataWithPagination[CustomerResponse]
