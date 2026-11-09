package address_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	address_enums "larsa-e-commerce-microservice-v2/internal/modules/address/enums"
)

type GeoLocationDTO struct {
	Type        string    `json:"type" validate:"required"`
	Coordinates []float64 `json:"coordinates" validate:"required,len=2"` // [longitude, latitude]
}

type AddressDTO struct {
	Name         types.Localizable[string] `json:"name" validate:"required"`
	Type         address_enums.AddressType `json:"type" validate:"required"`
	Country      string                    `json:"country" validate:"required"`
	City         string                    `json:"city" validate:"required"`
	Street       string                    `json:"street" validate:"required"`
	AddressLine1 string                    `json:"addressLine1" validate:"required"`
	Notes        *string                   `json:"notes"`
	POBox        *string                   `json:"poBox"`
	PostCode     *string                   `json:"postCode"`
	Phone        string                    `json:"phone" validate:"required"`
	Location     *GeoLocationDTO           `json:"location"`
	IsDefault    *bool                     `json:"isDefault"`
}
