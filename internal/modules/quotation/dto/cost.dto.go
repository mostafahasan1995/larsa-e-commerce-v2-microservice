package quotation_dto

type SetShippingCostDTO struct {
	FreightCharges *CostDto `json:"freightCharges"`
	Insurance      *CostDto `json:"insurance"`
}

type SetAgentFeeDTO struct {
	ExportDeclarationFees *CostDto `json:"exportDeclarationFees"`
	CustomsClearanceFee   *CostDto `json:"customsClearanceFee"`
	ImportDeclarationFees *CostDto `json:"importDeclarationFees"`
	DeliveryToBuyerFee    *CostDto `json:"deliveryToBuyerFee"`
}
