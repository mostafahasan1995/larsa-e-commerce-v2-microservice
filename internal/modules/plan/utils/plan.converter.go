package plan_utils

import (
	"context"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	plan_dto "larsa-e-commerce-microservice-v2/internal/modules/plan/dto"
	plan_models "larsa-e-commerce-microservice-v2/internal/modules/plan/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertCreatePlanDTOToModel(dto *plan_dto.CreatePlanDTO, ctx context.Context) (plan_models.Plan, *helpers.ApiError) {
	now := time.Now()

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return plan_models.Plan{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	plan := plan_models.Plan{
		BasePlan: plan_models.BasePlan{
			BaseModel: models.BaseModel{
				Id: primitive.NewObjectID(),
				MetaDataBaseModel: models.MetaDataBaseModel{
					CreatedAt: now,
					CreatedBy: &cfg.User.Id,
					UpdatedAt: now,
					UpdatedBy: &cfg.User.Id,
				},
			},
			Name:           dto.Name,
			Type:           dto.Type,
			TrailDays:      dto.TrailDays,
			IsDefault:      dto.IsDefault,
			IsActive:       dto.IsActive,
			ExpirationDays: dto.ExpirationDays,
			SuspendDays:    dto.SuspendDays,
			Description:    dto.Description,
			Image:          dto.Image,
		},
	}

	// Convert monthly prices
	for _, priceDTO := range dto.MonthlyPrices {
		price := plan_models.PlanPrice{
			BasePlanPrice: plan_models.BasePlanPrice{
				CurrencyId: priceDTO.CurrencyId,
				Price:      priceDTO.Price,
			},
		}
		plan.MonthlyPrices = append(plan.MonthlyPrices, price)
	}

	// Convert annual prices
	for _, priceDTO := range dto.AnnualPrices {
		price := plan_models.PlanPrice{
			BasePlanPrice: plan_models.BasePlanPrice{
				CurrencyId: priceDTO.CurrencyId,
				Price:      priceDTO.Price,
			},
		}
		plan.AnnualPrices = append(plan.AnnualPrices, price)
	}

	// Convert features
	for _, featureDTO := range dto.Features {
		feature := plan_models.PlanFeature{
			BasePlanFeature: plan_models.BasePlanFeature{
				FeatureId: featureDTO.FeatureId,
				Active:    featureDTO.Active,
				Value:     featureDTO.Value,
			},
		}
		plan.Features = append(plan.Features, feature)
	}

	// Convert customizations
	for _, customizationDTO := range dto.Customizations {
		customization := plan_models.PlanCustomization{
			BasePlanCustomization: plan_models.BasePlanCustomization{
				Type:     customizationDTO.Type,
				UsersIds: []primitive.ObjectID{},
			},
		}
		customization.UsersIds = customizationDTO.UsersIds
		plan.Customizations = append(plan.Customizations, customization)
	}

	return plan, nil
}

func ConvertUpdateDTOToModel(dto *plan_dto.UpdatePlanDTO, id primitive.ObjectID, ctx context.Context) (plan_models.Plan, *helpers.ApiError) {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return plan_models.Plan{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	now := time.Now()

	plan := plan_models.Plan{
		BasePlan: plan_models.BasePlan{
			BaseModel: models.BaseModel{
				Id: id,

				MetaDataBaseModel: models.MetaDataBaseModel{
					UpdatedAt: now,
					UpdatedBy: &cfg.User.Id,
				},
			},
		},
	}
	if dto.Name != nil {
		plan.Name = *dto.Name
	}
	if dto.Type != nil {
		plan.Type = *dto.Type
	}
	if dto.TrailDays != nil {
		plan.TrailDays = *dto.TrailDays
	}
	if dto.IsDefault != nil {
		plan.IsDefault = *dto.IsDefault
	}
	if dto.IsActive != nil {
		plan.IsActive = *dto.IsActive
	}
	if dto.ExpirationDays != nil {
		plan.ExpirationDays = *dto.ExpirationDays
	}
	if dto.SuspendDays != nil {
		plan.SuspendDays = *dto.SuspendDays
	}
	if dto.Description != nil {
		plan.Description = *dto.Description
	}
	if dto.Image != nil {
		plan.Image = *dto.Image
	}

	// Convert monthly prices
	if dto.MonthlyPrices != nil {
		plan.MonthlyPrices = nil
		for _, priceDTO := range dto.MonthlyPrices {
			price := plan_models.PlanPrice{}
			if priceDTO.CurrencyId != nil {
				price.CurrencyId = *priceDTO.CurrencyId
			}
			if priceDTO.Price != nil {
				price.Price = *priceDTO.Price
			}
			plan.MonthlyPrices = append(plan.MonthlyPrices, price)
		}
	}

	// Convert annual prices
	if dto.AnnualPrices != nil {
		plan.AnnualPrices = nil
		for _, priceDTO := range dto.AnnualPrices {
			price := plan_models.PlanPrice{}
			if priceDTO.CurrencyId != nil {
				price.CurrencyId = *priceDTO.CurrencyId
			}
			if priceDTO.Price != nil {
				price.Price = *priceDTO.Price
			}
			plan.AnnualPrices = append(plan.AnnualPrices, price)
		}
	}

	// Convert customizations
	if dto.Customizations != nil {
		plan.Customizations = nil
		for _, customizationDTO := range dto.Customizations {
			customization := plan_models.PlanCustomization{
				BasePlanCustomization: plan_models.BasePlanCustomization{
					Type:     *customizationDTO.Type,
					UsersIds: []primitive.ObjectID{},
				},
			}
			customization.UsersIds = customizationDTO.UsersIds
			plan.Customizations = append(plan.Customizations, customization)
		}
	}

	// Convert features
	if dto.Features != nil {
		plan.Features = nil
		for _, featureDTO := range dto.Features {
			feature := plan_models.PlanFeature{
				BasePlanFeature: plan_models.BasePlanFeature{
					FeatureId: *featureDTO.FeatureId,
					Active:    *featureDTO.Active,
					Value:     *featureDTO.Value,
				},
			}
			plan.Features = append(plan.Features, feature)
		}
	}

	return plan, nil
}
