package config_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	profile_models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseTicketType struct {
	models.MetaDataBaseModel `bson:",inline"`
	Code                     string               `bson:"code" json:"code"`
	Name                     string               `bson:"name" json:"name"`
	ExpectedSolveTime        float64              `bson:"expectedSolveTime" json:"expectedSolveTime"` // in minutes
	UpdateNotification       bool                 `bson:"updateNotification" json:"updateNotification"`
	InternalNotification     bool                 `bson:"internalNotification" json:"internalNotification"`
	TeamNotifyIds            []primitive.ObjectID `bson:"teamNotifyIds" json:"teamNotifyIds"`
	IsActive                 *bool                `bson:"isActive,omitempty" json:"isActive,omitempty"`
}

type TicketType struct {
	BaseTicketType `bson:",inline"`
}

type TicketTypeResponse struct {
	BaseTicketType `bson:",inline"`
	TeamNotify     []profile_models.ProfileResponse `bson:"teamNotify" json:"teamNotify"`
}

type TicketTypesWithPagination = types.DataWithPagination[TicketTypeResponse]
