package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	Id                primitive.ObjectID `bson:"_id" json:"id"`
	MetaDataBaseModel `bson:",inline"`
}

type MetaDataBaseModel struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	CreatedAt time.Time           `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy *primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt time.Time           `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy *primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	DeletedAt *time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	DeletedBy *primitive.ObjectID `bson:"deletedBy,omitempty" json:"deletedBy,omitempty"`
}

type OptionalMetaDataBaseModel struct {
	Id *primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	CreatedAt *time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	CreatedBy *primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedAt *time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	UpdatedBy *primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	DeletedAt *time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	DeletedBy *primitive.ObjectID `bson:"deletedBy,omitempty" json:"deletedBy,omitempty"`
}
