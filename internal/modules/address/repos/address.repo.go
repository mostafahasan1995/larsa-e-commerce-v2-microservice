package address_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	address_filters "larsa-e-commerce-microservice-v2/internal/modules/address/filters"
	address_models "larsa-e-commerce-microservice-v2/internal/modules/address/models"
	address_services "larsa-e-commerce-microservice-v2/internal/modules/address/services"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddressRepository implements the repository interface
type AddressRepository struct {
	*base_repos.BaseRepo[address_models.AddressModel, address_models.AddressResponse]
	db         *mongo.Client
	collection enums.CollectionName
}

// NewAddressRepository creates a new instance of AddressRepository
func NewAddressRepository(i *do.Injector) (address_services.AddressRepository, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &AddressRepository{
		BaseRepo:   base_repos.NewBaseRepo[address_models.AddressModel, address_models.AddressResponse](db, enums.ADDRESS),
		db:         db,
		collection: enums.ADDRESS,
	}, nil
}

// Create creates a new address
func (r *AddressRepository) Create(ctx context.Context, model address_models.AddressModel) (*address_models.AddressModel, *helpers.ApiError) {
	// Check if this is the default address and set others to non-default if needed
	if model.IsDefault != nil && *model.IsDefault && model.UserId != nil {
		err := r.updateDefaultAddressStatus(ctx, *model.UserId, false)
		if err != nil {
			return nil, err
		}
	}

	// Ensure ID is set
	if model.Id == primitive.NilObjectID {
		model.Id = primitive.NewObjectID()
	}

	// Set timestamps
	now := time.Now()
	model.CreatedAt = now
	model.UpdatedAt = now

	_, err := r.BaseRepo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// GetById retrieves an address by ID
func (r *AddressRepository) GetById(ctx context.Context, id primitive.ObjectID) (*address_models.AddressModel, *helpers.ApiError) {
	// Build pipeline for getting a single address
	matchStage := bson.M{}
	filter, _ := address_filters.NewAddressFilters("")
	pipeline := filter.BuildPipeline(matchStage, ctx, false)

	response, err := r.BaseRepo.GetById(ctx, id, pipeline)
	if err != nil {
		return nil, err
	}

	return &response.AddressModel, nil
}

// UpdateById updates an address by ID
func (r *AddressRepository) UpdateById(ctx context.Context, id primitive.ObjectID, model address_models.AddressModel) (*address_models.AddressModel, *helpers.ApiError) {
	if model.IsDefault != nil && *model.IsDefault && model.UserId != nil {
		err := r.updateDefaultAddressStatus(ctx, *model.UserId, false)
		if err != nil {
			return nil, err
		}
	}

	model.UpdatedAt = time.Now()

	response, err := r.BaseRepo.UpdateById(ctx, id, model, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return nil, err
	}

	return &response.AddressModel, nil
}

// DeleteById soft deletes an address by ID
func (r *AddressRepository) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return r.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
}

// BulkDelete soft deletes multiple addresses
func (r *AddressRepository) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return r.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
}

// SetDefaultAddress sets an address as default for a user
func (r *AddressRepository) SetDefaultAddress(ctx context.Context, id, userId primitive.ObjectID) *helpers.ApiError {
	updateErr := r.updateDefaultAddressStatus(ctx, userId, false)
	if updateErr != nil {
		return updateErr
	}

	cfg, _ := util.GetReqAppCfg(ctx)

	// Then set the specified address as default
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"isDefault": true,
			"updatedAt": now,
			"updatedBy": cfg.User.Id,
		},
	}

	collection, _ := r.BaseRepo.GetCollection(ctx)
	result, mongoErr := collection.UpdateOne(
		ctx,
		bson.M{"_id": id, "userId": userId},
		update,
	)

	if mongoErr != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    "Error setting default address: " + mongoErr.Error(),
		}
	}

	if result.MatchedCount == 0 {
		return &helpers.ApiError{
			StatusCode: 404,
			Message:    "Address not found",
		}
	}

	return nil
}

// Helper function to update the default status of all addresses for a user
func (r *AddressRepository) updateDefaultAddressStatus(ctx context.Context, userId primitive.ObjectID, isDefault bool) *helpers.ApiError {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get user context",
		}
	}

	update := bson.M{
		"$set": bson.M{
			"isDefault": isDefault,
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	collection, _ := r.BaseRepo.GetCollection(ctx)
	_, updateErr := collection.UpdateMany(
		ctx,
		bson.M{"userId": userId, "deletedAt": nil},
		update,
	)

	if updateErr != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    "Error updating default address status: " + updateErr.Error(),
		}
	}

	return nil
}

// GetAll retrieves all addresses
func (repo *AddressRepository) GetAll(ctx context.Context, filter *address_filters.AddressFilters) ([]address_models.AddressResponse, *helpers.ApiError) {
	// Build pipeline with filters
	matchStage := bson.M{}
	//TODO: I change the extra operations to false. we need to check if it's working
	pipeline := filter.BuildPipeline(matchStage, ctx, false)

	// Get all addresses
	return repo.BaseRepo.GetAll(ctx, pipeline)
}

// GetWithPagination retrieves addresses with pagination
func (repo *AddressRepository) GetWithPagination(ctx context.Context, skip, limit int64, filter *address_filters.AddressFilters) (*address_models.AddressesWithPagination, *helpers.ApiError) {
	// Build pipeline with filters
	matchStage := bson.M{}
	pipeline := filter.BuildPipeline(matchStage, ctx, true)

	// Get addresses with pagination
	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

// MigrateAddresses migrates all addresses from old user to new user without changing the address IDs
func (r *AddressRepository) MigrateAddresses(ctx context.Context, oldUserId, newUserId primitive.ObjectID) *helpers.ApiError {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get user context",
		}
	}

	// Update all addresses belonging to the old user to the new user
	update := bson.M{
		"$set": bson.M{
			"userId":    newUserId,
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	collection, _ := r.BaseRepo.GetCollection(ctx)
	_, updateErr := collection.UpdateMany(
		ctx,
		bson.M{"userId": oldUserId, "deletedAt": nil},
		update,
	)

	if updateErr != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    "Error migrating addresses: " + updateErr.Error(),
		}
	}

	return nil
}
