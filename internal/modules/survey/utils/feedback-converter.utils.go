package survey_utils

import (
	"context"
	"crypto/rand"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	survey_dtos "larsa-e-commerce-microservice-v2/internal/modules/survey/dtos"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"
	"math/big"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertCreateFeedbackDTOToModel(dto *survey_dtos.CreateFeedbackDTO, ctx context.Context) (survey_models.FeedbackModel, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return survey_models.FeedbackModel{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	userId := primitive.NilObjectID

	if cfg.User != nil {
		userId = cfg.User.Id
	}

	additionalOptions := types.AdditionalOptions{
		UserId: userId,
		Mode:   types.CREATE,
	}
	email := ""
	if dto.Email != nil {
		email = *dto.Email
	}

	return survey_models.FeedbackModel{
		Phone:             dto.Phone,
		ClientName:        dto.ClientName,
		Email:             email,
		Birthday:          dto.Birthday,
		WaiterId:          dto.WaiterId,
		BranchId:          dto.BranchId,
		VendorId:          dto.VendorId,
		Points:            dto.Points,
		Answers:           dto.Answers,
		Notes:             dto.Notes,
		WaiterRate:        dto.WaiterRate,
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
	}, nil
}

// GenerateCouponCode generates a unique 8-character coupon code
func GenerateCouponCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}

	return string(b)
}

func ConvertUpdateFeedbackDTOToModel(dto *survey_dtos.UpdateFeedbackDTO, id primitive.ObjectID, ctx context.Context) (survey_models.FeedbackModel, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return survey_models.FeedbackModel{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	additionalOptions := types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	}

	return survey_models.FeedbackModel{
		Phone:             dto.Phone,
		ClientName:        dto.ClientName,
		Email:             dto.Email,
		Birthday:          dto.Birthday,
		WaiterId:          dto.WaiterId,
		BranchId:          dto.BranchId,
		VendorId:          dto.VendorId,
		Points:            dto.Points,
		Answers:           dto.Answers,
		Notes:             dto.Notes,
		WaiterRate:        dto.WaiterRate,
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
	}, nil
}
