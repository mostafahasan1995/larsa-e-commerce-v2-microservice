package config_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTicketTypeDto struct {
	Code                 string               `json:"code" validate:"required"`
	Name                 string               `json:"name" validate:"required"`
	ExpectedSolveTime    float64              `json:"expectedSolveTime" validate:"required"` // in minutes
	UpdateNotification   bool                 `json:"updateNotification"`
	InternalNotification bool                 `json:"internalNotification"`
	TeamNotifyIds        []primitive.ObjectID `json:"teamNotifyIds"`
	IsActive             *bool                `json:"isActive,omitempty"`
}

func (dto *CreateTicketTypeDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	// Validate that expected solve time is positive
	if dto.ExpectedSolveTime <= 0 {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Expected solve time must be greater than 0",
		}
	}

	// If team notify IDs are provided, validate that they exist
	if len(dto.TeamNotifyIds) > 0 {
		profileRepo := base_repos.NewBaseRepo[any, any](args.DB, enums.PROFILE)
		for _, teamId := range dto.TeamNotifyIds {
			_, err := profileRepo.GetById(args.Context, teamId, nil)
			if err != nil {
				return nil, &helpers.ApiError{
					StatusCode: 400,
					Message:    "Invalid team ID in teamNotifyIds",
				}
			}
		}
	}

	return nil, nil
}
