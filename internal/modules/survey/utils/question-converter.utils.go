package survey_utils

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	survey_dtos "larsa-e-commerce-microservice-v2/internal/modules/survey/dtos"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertCreateQuestionDTOToModel(dto *survey_dtos.CreateQuestionDTO, ctx context.Context) (survey_models.QuestionModel, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return survey_models.QuestionModel{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	additionalOptions := types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	}

	return survey_models.QuestionModel{
		Type:              dto.Type,
		Order:             dto.Order,
		Question:          dto.Question,
		Options:           dto.Options,
		BranchIds:         dto.BranchIds,
		AllBranches:       dto.AllBranches,
		VendorId:          dto.VendorId,
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
	}, nil
}

func ConvertUpdateQuestionDTOToModel(dto *survey_dtos.UpdateQuestionDTO, id primitive.ObjectID, ctx context.Context) (survey_models.QuestionModel, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return survey_models.QuestionModel{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	additionalOptions := types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	}

	return survey_models.QuestionModel{
		Type:              dto.Type,
		Order:             dto.Order,
		Question:          dto.Question,
		Options:           dto.Options,
		BranchIds:         dto.BranchIds,
		AllBranches:       dto.AllBranches,
		VendorId:          dto.VendorId,
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
	}, nil
}
