package payment_dtos

type CreatePayoutDto struct {
	Amount            float64 `json:"amount" validate:"required,min=0"`
	BankName          string  `json:"bankName" validate:"required"`
	BankAccountNumber string  `json:"bankAccountNumber" validate:"required"`
	AccountHolderName string  `json:"accountHolderName" validate:"required"`
	BankBranchAddress string  `json:"bankBranchAddress" validate:"required"`
	SwiftCode         string  `json:"swiftCode" validate:"required"`
	TransferReason    string  `json:"transferReason" validate:"required"`
	Notes             string  `json:"notes"`
}
