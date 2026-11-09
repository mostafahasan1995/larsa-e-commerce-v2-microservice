package unit_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	unit_enums "larsa-e-commerce-microservice-v2/internal/modules/unit/enums"
	unit_models "larsa-e-commerce-microservice-v2/internal/modules/unit/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type CreateUnitDto struct {
	Name     types.Localizable[string] `json:"name" validate:"required"`
	Type     unit_enums.UnitType       `json:"type" validate:"required,oneof=bigger_than_reference reference smaller_than_reference"`
	Category unit_enums.UnitCategory   `json:"category" validate:"required,oneof=length unit weight volume working_time area"`
	Rate     *float32                  `json:"rate" validate:"omitempty,min=0"`
	Active   bool                      `json:"active"`
}

func (dto *CreateUnitDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	unitRepo := base_repos.NewBaseRepo[unit_models.Unit, unit_models.UnitResponse](args.DB, enums.UNIT)
	if dto.Type == unit_enums.Reference && dto.Rate != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Rate cannot be set when type is reference",
		}
	}

	if dto.Type != unit_enums.Reference && dto.Rate == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Rate is required when type is not reference",
		}
	}

	if dto.Type == unit_enums.SmallerThanReference && *dto.Rate > 1 {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Rate must be smaller than 1 when type is smaller_than_reference",
		}
	}

	if dto.Type == unit_enums.BiggerThanReference && *dto.Rate <= 1 {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Rate must be bigger than 1 when type is bigger_than_reference",
		}
	}

	softDeletedFilter := base_utils.AddSoftDeletedFilter(bson.M{})
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"type":     unit_enums.Reference,
				"category": dto.Category,
			},
		},
		{
			"$match": softDeletedFilter,
		},
	}
	storedUnit, _ := unitRepo.FindOne(args.Context, pipeline)

	/*
	*	if the unit is reference and a reference unit already exists for this category, return an error
	 */
	if dto.Type == unit_enums.Reference && storedUnit != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "A reference unit already exists for this category",
		}
	}

	/*
	*	if the unit is not reference and no reference unit exists for this category, return an error
	 */
	if dto.Type != unit_enums.Reference && storedUnit == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "no reference unit exists for this category",
		}
	}

	return nil, nil
}
