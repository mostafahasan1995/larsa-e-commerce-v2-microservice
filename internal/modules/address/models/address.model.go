package address_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	address_enums "larsa-e-commerce-microservice-v2/internal/modules/address/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressModel struct {
	models.MetaDataBaseModel `bson:",inline" json:",inline"`

	Name  *types.Localizable[string] `bson:"name,omitempty" json:"name,omitempty"`
	Guest *bool                      `bson:"guest,omitempty" json:"guest,omitempty"`
	Type  *address_enums.AddressType `bson:"type,omitempty" json:"type,omitempty"`

	Country      *string      `bson:"country,omitempty" json:"country,omitempty"`
	City         *string      `bson:"city,omitempty" json:"city,omitempty"`
	Street       *string      `bson:"street,omitempty" json:"street,omitempty"`
	AddressLine1 *string      `bson:"addressLine1,omitempty" json:"addressLine1,omitempty"`
	Notes        *string      `bson:"notes,omitempty" json:"notes,omitempty"`
	POBox        *string      `bson:"poBox,omitempty" json:"poBox,omitempty"`
	PostCode     *string      `bson:"postCode,omitempty" json:"postCode,omitempty"`
	Phone        *string      `bson:"phone,omitempty" json:"phone,omitempty"`
	Location     *GeoLocation `bson:"location,omitempty" json:"location,omitempty"`

	IsDefault *bool               `bson:"isDefault,omitempty" json:"isDefault,omitempty"`
	UserId    *primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
}

type GeoLocation struct {
	Type        string    `bson:"type" json:"type" validate:"required,oneof=Point"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [longitude, latitude]
}

type AddressResponse struct {
	AddressModel `bson:",inline" json:",inline"`
}

type AddressesWithPagination = types.DataWithPagination[AddressResponse]
