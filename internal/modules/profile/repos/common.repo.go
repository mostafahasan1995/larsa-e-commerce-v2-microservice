package profile_repos

import (
	"context"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	profile_filters "larsa-e-commerce-microservice-v2/internal/modules/profile/filters"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetById retrieves a profile by its ID
func (repo *ProfileRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProfileResponse, *helpers.ApiError) {
	collection, _ := repo.BaseRepo.GetCollection(ctx)
	filter := bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}
	var profile models.ProfileResponse
	err := collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &helpers.ApiError{StatusCode: 404, Message: "profile not found"}
		}
		return nil, helpers.DefaultError(err)
	}
	return &profile, nil
}

// GetWithPagination retrieves profiles with pagination
func (repo *ProfileRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProfilesWithPagination, *helpers.ApiError) {
	profileFilters, err := profile_filters.NewProfileFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "invalid filter format: " + err.Error()}
	}

	matchStage := bson.M{}
	pipeline := profileFilters.BuildPipeline(matchStage, ctx, true)
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	collection, _ := repo.BaseRepo.GetCollection(ctx)
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}
	defer cursor.Close(ctx)

	var profiles []models.ProfileResponse
	if err = cursor.All(ctx, &profiles); err != nil {
		return nil, helpers.DefaultError(err)
	}

	countPipeline := profileFilters.BuildPipeline(matchStage, ctx, false)
	countPipeline = append(countPipeline, bson.M{"$count": "total"})

	countCursor, err := collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}
	defer countCursor.Close(ctx)

	var countResult []bson.M
	if err = countCursor.All(ctx, &countResult); err != nil {
		return nil, helpers.DefaultError(err)
	}

	var totalCount int64 = 0
	if len(countResult) > 0 {
		if total, ok := countResult[0]["total"].(int32); ok {
			totalCount = int64(total)
		}
	}

	totalPages := float64(totalCount) / float64(limit)
	if totalCount%limit > 0 {
		totalPages++
	}

	result := &models.ProfilesWithPagination{
		Data:       profiles,
		Pagination: types.Pagination{TotalPages: totalPages, PerPage: limit, TotalCount: totalCount},
	}

	return result, nil
}

// GetProfileByUserId retrieves a profile by user ID
func (repo *ProfileRepo) GetProfileByUserId(ctx context.Context, userId primitive.ObjectID) (*models.ProfileResponse, *helpers.ApiError) {
	return repo.GetById(ctx, userId, "")
}

// DeleteById deletes a profile by its ID
func (repo *ProfileRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	profileFilters, _ := profile_filters.NewProfileFilters("")
	matchStage := bson.M{"deletedAt": bson.M{"$exists": false}}
	pipeline := profileFilters.BuildPipeline(matchStage, ctx, false)

	_, getErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if getErr != nil {
		return getErr
	}

	now := time.Now()
	collection, _ := repo.BaseRepo.GetCollection(ctx)
	_, deleteErr := collection.UpdateOne(
		ctx,
		bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"deletedAt": now, "deletedBy": cfg.User.Id, "updatedAt": now, "updatedBy": cfg.User.Id}},
	)
	if deleteErr != nil {
		return helpers.DefaultError(deleteErr)
	}
	return nil
}

// BulkDelete deletes multiple profiles by their IDs
func (repo *ProfileRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	profileFilters, _ := profile_filters.NewProfileFilters("")

	if cfg.User.Type != enums.SuperAdmin {
		for _, id := range ids {
			matchStage := bson.M{"deletedAt": bson.M{"$exists": false}}
			pipeline := profileFilters.BuildPipeline(matchStage, ctx, false)
			_, getErr := repo.BaseRepo.GetById(ctx, id, pipeline)
			if getErr != nil {
				return &helpers.ApiError{StatusCode: 403, Message: "not authorized to delete profile with id: " + id.Hex()}
			}
		}
	}

	now := time.Now()
	collection, _ := repo.BaseRepo.GetCollection(ctx)
	_, deleteErr := collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": ids}, "deletedAt": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"deletedAt": now, "deletedBy": cfg.User.Id, "updatedAt": now, "updatedBy": cfg.User.Id}},
	)
	if deleteErr != nil {
		return helpers.DefaultError(deleteErr)
	}
	return nil
}

// GetMe retrieves the authenticated user's profile
func (repo *ProfileRepo) GetMe(ctx context.Context) (*models.ProfileResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	if cfg == nil || cfg.User == nil || cfg.User.Id.IsZero() {
		return nil, &helpers.ApiError{StatusCode: 401, Message: "user not authenticated"}
	}
	userId := cfg.User.Id
	profile, err := repo.GetProfileByUserId(ctx, userId)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 404, Message: "profile not found"}
	}
	return profile, nil
}

// GetProfileByFingerprintId retrieves a profile by fingerprint ID
func (repo *ProfileRepo) GetProfileByFingerprintId(ctx context.Context, fingerprintId string) (*models.ProfileResponse, *helpers.ApiError) {
	collection, _ := repo.BaseRepo.GetCollection(ctx)
	filter := bson.M{"fingerprintId": fingerprintId, "deletedAt": bson.M{"$exists": false}}
	var profile models.ProfileResponse
	err := collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, helpers.DefaultError(err)
	}
	return &profile, nil
}
