package appointment_models

import (
	appointment_enums "larsa-e-commerce-microservice-v2/internal/modules/appointment/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseTranslation struct {
	TranslateType appointment_enums.TranslateType `bson:"translateType" json:"translateType"`
}

type Translation struct {
	BaseTranslation `bson:",inline"`

	TranslateFromId *primitive.ObjectID `bson:"translateFromId" json:"translateFromId,omitempty"`
	TranslateToId   *primitive.ObjectID `bson:"translateToId" json:"translateToId,omitempty"`
	InvoiceId       *primitive.ObjectID `bson:"invoiceId" json:"invoiceId,omitempty"`
}

type TranslationResponse struct {
	BaseTranslation `bson:",inline"`

	TranslateFrom string `bson:"translateFrom" json:"translateFrom,omitempty"`
	TranslateTo   string `bson:"translateTo" json:"translateTo,omitempty"`
	Invoice       string `bson:"invoice" json:"invoice,omitempty"`
}
