package profile_repos

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	profile_utils "larsa-e-commerce-microservice-v2/internal/modules/profile/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateBuyerProfile creates a new buyer profile
func (repo *ProfileRepo) CreateBuyerProfile(ctx context.Context, dto *dtos.CreateBuyerProfileDto) (*models.ProfileResponse, *helpers.ApiError) {
	if dto.BaseProfileDto.UserId != nil {
		userId, err := primitive.ObjectIDFromHex(*dto.BaseProfileDto.UserId)
		if err != nil {
			return nil, helpers.DefaultError(err)
		}
		existingProfile, _ := repo.GetById(ctx, userId, "")
		if existingProfile != nil {
			return nil, &helpers.ApiError{StatusCode: 409, Message: "profile already exists for this user"}
		}
	}

	if dto.Password == nil {
		password := "123456"
		dto.Password = &password
	}
	if dto.FirstName == nil {
		firstName := "Buyer"
		dto.FirstName = &firstName
	}
	if dto.LastName == nil {
		lastName := ""
		dto.LastName = &lastName
	}

	userId, err := repo.usersgw.AddUser(ctx, &types.PostUserData{
		Email:     *dto.Email,
		Password:  *dto.Password,
		FirstName: *dto.FirstName,
		LastName:  *dto.LastName,
	})
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "failed to create user: " + err.Error()}
	}

	userIdHex := userId.Hex()
	dto.BaseProfileDto.UserId = &userIdHex

	return repo.CreateBuyerProfileFromUser(ctx, dto, userId)
}

// UpdateBuyerProfile updates an existing buyer profile
func (repo *ProfileRepo) UpdateBuyerProfile(ctx context.Context, id primitive.ObjectID, dto *dtos.CreateBuyerProfileDto) (*models.ProfileResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	matchStage := bson.M{}

	profile := profile_utils.ConvertToBuyerProfileModel(dto, cfg.User.Id)
	profile.Id = id

	commonProfile := models.CommonProfile{BaseProfile: profile.BaseProfile}

	return repo.BaseRepo.UpdateById(ctx, id, commonProfile, matchStage)
}

func (repo *ProfileRepo) CreateBuyerProfileFromUser(ctx context.Context, dto *dtos.CreateBuyerProfileDto, userId primitive.ObjectID) (*models.ProfileResponse, *helpers.ApiError) {
	profile := profile_utils.ConvertToBuyerProfileModel(dto, userId)
	collection, _ := repo.BaseRepo.GetCollection(ctx)
	_, err := collection.InsertOne(ctx, profile)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	result, getErr := repo.GetById(ctx, profile.Id, "")
	if getErr != nil {
		return nil, getErr
	}

	return result, nil
}
