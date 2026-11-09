package address_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	address_dtos "larsa-e-commerce-microservice-v2/internal/modules/address/dtos"
	address_filters "larsa-e-commerce-microservice-v2/internal/modules/address/filters"
	address_models "larsa-e-commerce-microservice-v2/internal/modules/address/models"
	address_utils "larsa-e-commerce-microservice-v2/internal/modules/address/utils"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressRepository interface {
	Create(ctx context.Context, model address_models.AddressModel) (*address_models.AddressModel, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID) (*address_models.AddressModel, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, model address_models.AddressModel) (*address_models.AddressModel, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
	SetDefaultAddress(ctx context.Context, id, userId primitive.ObjectID) *helpers.ApiError
	GetAll(ctx context.Context, filter *address_filters.AddressFilters) ([]address_models.AddressResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter *address_filters.AddressFilters) (*address_models.AddressesWithPagination, *helpers.ApiError)
	MigrateAddresses(ctx context.Context, oldUserId, newUserId primitive.ObjectID) *helpers.ApiError
}

type AddressService interface {
	Create(ctx context.Context, dto address_dtos.AddressDTO) (*address_models.AddressResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID) (*address_models.AddressResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, dto address_dtos.AddressDTO) (*address_models.AddressResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
	SetDefaultAddress(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	GetAll(ctx context.Context, query string) ([]address_models.AddressResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*address_models.AddressesWithPagination, *helpers.ApiError)
	MigrateAddresses(ctx context.Context, oldUserId, newUserId primitive.ObjectID) *helpers.ApiError
}

// AddressServiceImpl implements the AddressService interface
type AddressServiceImpl struct {
	addressRepo AddressRepository
}

// NewAddressService creates a new AddressService implementation
func NewAddressService(i *do.Injector) (AddressService, error) {
	return &AddressServiceImpl{
		addressRepo: do.MustInvoke[AddressRepository](i),
	}, nil
}

// convertModelToResponseDTO converts a model to a response DTO
func (s *AddressServiceImpl) convertModelToResponseDTO(model *address_models.AddressModel) address_models.AddressResponse {
	// Create a response using the existing model
	response := address_models.AddressResponse{
		AddressModel: *model,
	}

	return response
}

// Create creates a new address
func (s *AddressServiceImpl) Create(ctx context.Context, dto address_dtos.AddressDTO) (*address_models.AddressResponse, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get user context",
		}
	}
	model := address_utils.ConvertAddressDTOToModel(dto, cfg.User.Id)

	// Call repository to create the address
	result, apiErr := s.addressRepo.Create(ctx, *model)
	if apiErr != nil {
		return nil, apiErr
	}

	// Convert model to response DTO
	response := s.convertModelToResponseDTO(result)
	return &response, nil
}

// GetById gets an address by ID
func (s *AddressServiceImpl) GetById(ctx context.Context, id primitive.ObjectID) (*address_models.AddressResponse, *helpers.ApiError) {
	// Call repository to get the address
	result, apiErr := s.addressRepo.GetById(ctx, id)
	if apiErr != nil {
		return nil, apiErr
	}

	// Convert model to response DTO
	response := s.convertModelToResponseDTO(result)
	return &response, nil
}

// UpdateById updates an address by ID
func (s *AddressServiceImpl) UpdateById(ctx context.Context, id primitive.ObjectID, dto address_dtos.AddressDTO) (*address_models.AddressResponse, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get user context",
		}
	}
	model := address_utils.ConvertAddressDTOToModel(dto, cfg.User.Id)
	model.UpdatedAt = time.Now()
	model.Id = id

	// Call repository to update the address
	result, apiErr := s.addressRepo.UpdateById(ctx, id, *model)
	if apiErr != nil {
		return nil, apiErr
	}

	// Convert model to response DTO
	response := s.convertModelToResponseDTO(result)
	return &response, nil
}

// DeleteById deletes an address by ID
func (s *AddressServiceImpl) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.addressRepo.DeleteById(ctx, id)
}

// BulkDelete deletes multiple addresses
func (s *AddressServiceImpl) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return s.addressRepo.BulkDelete(ctx, ids)
}

// SetDefaultAddress sets an address as default for a user
func (s *AddressServiceImpl) SetDefaultAddress(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get user context",
		}
	}

	return s.addressRepo.SetDefaultAddress(ctx, id, cfg.User.Id)
}

// GetAll gets all addresses
func (s *AddressServiceImpl) GetAll(ctx context.Context, query string) ([]address_models.AddressResponse, *helpers.ApiError) {
	// Create filters from query parameters
	//userId
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "invalid filter format: " + err.Error(),
		}
	}
	addressFilters, err := address_filters.NewAddressFilters(query)

	addressFilters.UserId = &cfg.User.Id
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "invalid filter format: " + err.Error(),
		}
	}

	return s.addressRepo.GetAll(ctx, addressFilters)
}

// GetWithPagination gets addresses with pagination
func (s *AddressServiceImpl) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*address_models.AddressesWithPagination, *helpers.ApiError) {
	//userId
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "invalid filter format: " + err.Error(),
		}
	}
	addressFilters, err := address_filters.NewAddressFilters(query)
	addressFilters.UserId = &cfg.User.Id
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "invalid filter format: " + err.Error(),
		}
	}

	return s.addressRepo.GetWithPagination(ctx, skip, limit, addressFilters)
}

// MigrateAddresses migrates all addresses from old user to new user
func (s *AddressServiceImpl) MigrateAddresses(ctx context.Context, oldUserId, newUserId primitive.ObjectID) *helpers.ApiError {
	return s.addressRepo.MigrateAddresses(ctx, oldUserId, newUserId)
}
