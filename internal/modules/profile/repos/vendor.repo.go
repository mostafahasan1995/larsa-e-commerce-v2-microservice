package profile_repos

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	profile_utils "larsa-e-commerce-microservice-v2/internal/modules/profile/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateVendorProfile creates a new vendor profile
func (repo *ProfileRepo) CreateVendorProfile(ctx context.Context, dto *dtos.CreateVendorProfileDto) (*models.ProfileResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	userId, err := primitive.ObjectIDFromHex(*dto.BaseProfileDto.UserId)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}
	existingProfile, _ := repo.GetById(ctx, userId, "")
	if existingProfile != nil {
		return nil, &helpers.ApiError{StatusCode: 409, Message: "profile already exists for this user"}
	}

	profile := profile_utils.ConvertToVendorProfileModel(dto, cfg.User.Id)

	collection, _ := repo.BaseRepo.GetCollection(ctx)
	_, err = collection.InsertOne(ctx, profile)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	result, getErr := repo.GetById(ctx, profile.Id, "")
	if getErr != nil {
		return nil, getErr
	}
	return result, nil
}

// UpdateVendorProfile updates an existing vendor profile
func (repo *ProfileRepo) UpdateVendorProfile(ctx context.Context, id primitive.ObjectID, dto *dtos.CreateVendorProfileDto) (*models.ProfileResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	matchStage := bson.M{}

	profile := profile_utils.ConvertToVendorProfileModel(dto, cfg.User.Id)
	profile.Id = id

	commonProfile := models.CommonProfile{
		BaseProfile:                profile.BaseProfile,
		RegistrationCertifications: profile.RegistrationCertifications,
		QualityCertifications:      profile.QualityCertifications,
		Licenses:                   profile.Licenses,
	}

	return repo.BaseRepo.UpdateById(ctx, id, commonProfile, matchStage)
}
