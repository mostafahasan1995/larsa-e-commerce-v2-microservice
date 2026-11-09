package product_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Option struct {
	Label string `bson:"label" json:"label"`
	Value string `bson:"value" json:"value"`
}

type BaseAttribute struct {
	Name                          types.Localizable[string]        `bson:"name" json:"name"`
	Options                       types.LocalizableArray[Option]   `bson:"options" json:"options"`
	Code                          string                           `bson:"code" json:"code"`
	Type                          product_enums.AttributeType      `bson:"type" json:"type"`
	Tag                           product_enums.AttributeTag       `bson:"tag" json:"tag"`
	IsFiltrable                   bool                             `bson:"isFiltrable" json:"isFiltrable"`
	Required                      bool                             `bson:"required" json:"required"`
	Multi                         bool                             `bson:"multi" json:"multi"`
	IsComparable                  bool                             `bson:"isComparable" json:"isComparable"`
	IsLocalizable                 bool                             `bson:"isLocalizable" json:"isLocalizable"`
	RelatedTo                     *[]string                        `bson:"relatedTo,omitempty" json:"relatedTo,omitempty"`
	Visibility                    bool                             `bson:"visibility" json:"visibility"`
	Sortable                      bool                             `bson:"sortable" json:"sortable"`
	SelectionStyle                *string                          `bson:"selectionStyle" json:"selectionStyle"`
	BelongsTo                     []product_enums.BelongTo         `bson:"belongsTo,omitempty" json:"belongsTo,omitempty"`
	GroupId                       *primitive.ObjectID              `bson:"groupId,omitempty" json:"groupId,omitempty"`
	RefField                      *product_enums.AttributeRefField `bson:"refField,omitempty" json:"refField,omitempty"`
	VendorId                      *primitive.ObjectID              `bson:"vendorId,omitempty" json:"vendorId,omitempty"`
	Configurable                  bool                             `bson:"configurable" json:"configurable"`
	base_models.MetaDataBaseModel `bson:",inline"`
}

type Attribute struct {
	BaseAttribute `bson:",inline"`
}

type AttributeResponse struct {
	BaseAttribute `bson:",inline"`
	Group         *GroupResponse `bson:"group" json:"group"`
}

type AttributesWithPagination = types.DataWithPagination[AttributeResponse]
