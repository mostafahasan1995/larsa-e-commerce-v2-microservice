package profile_services

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/gateway"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProfileRepo defines the repository interface for profiles
type ProfileRepo interface {
	// Vendor profile methods
	CreateVendorProfile(ctx context.Context, dto *dtos.CreateVendorProfileDto) (*models.ProfileResponse, *helpers.ApiError)
	UpdateVendorProfile(ctx context.Context, id primitive.ObjectID, dto *dtos.CreateVendorProfileDto) (*models.ProfileResponse, *helpers.ApiError)

	// Buyer profile methods
	CreateBuyerProfile(ctx context.Context, dto *dtos.CreateBuyerProfileDto) (*models.ProfileResponse, *helpers.ApiError)
	UpdateBuyerProfile(ctx context.Context, id primitive.ObjectID, dto *dtos.CreateBuyerProfileDto) (*models.ProfileResponse, *helpers.ApiError)
	CreateBuyerProfileFromUser(ctx context.Context, dto *dtos.CreateBuyerProfileDto, userId primitive.ObjectID) (*models.ProfileResponse, *helpers.ApiError)

	// Guest profile methods
	CreateGuestProfile(ctx context.Context, dto *dtos.CreateGuestProfileDto, userId primitive.ObjectID) (*models.ProfileResponse, *helpers.ApiError)
	GetProfileByFingerprintId(ctx context.Context, fingerprintId string) (*models.ProfileResponse, *helpers.ApiError)
	DeleteAllGuestProfiles(ctx context.Context) (int64, *helpers.ApiError)

	// User profile
	GetMe(ctx context.Context) (*models.ProfileResponse, *helpers.ApiError)

	// Common methods
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProfileResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProfilesWithPagination, *helpers.ApiError)
	GetProfileByUserId(ctx context.Context, userId primitive.ObjectID) (*models.ProfileResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
	GetAll(ctx context.Context, pipeline []bson.M) ([]models.ProfileResponse, *helpers.ApiError)
}

// ProfileService defines the service interface for profiles
type ProfileService interface {
	// Vendor profile methods
	CreateVendorProfile(ctx context.Context, dto *dtos.CreateVendorProfileDto) (*models.ProfileResponse, *helpers.ApiError)
	UpdateVendorProfile(ctx context.Context, id primitive.ObjectID, dto *dtos.CreateVendorProfileDto) (*models.ProfileResponse, *helpers.ApiError)

	// Buyer profile methods
	CreateBuyerProfile(ctx context.Context, dto *dtos.CreateBuyerProfileDto) (*models.ProfileResponse, *helpers.ApiError)
	UpdateBuyerProfile(ctx context.Context, id primitive.ObjectID, dto *dtos.CreateBuyerProfileDto) (*models.ProfileResponse, *helpers.ApiError)

	// Guest profile methods
	CreateGuestProfile(ctx *context.Context, dto *dtos.CreateGuestProfileDto) (*models.ProfileResponse, *helpers.ApiError)
	GetProfileByFingerprintId(ctx context.Context, fingerprintId string) (*models.ProfileResponse, *helpers.ApiError)
	DeleteAllGuestProfiles(ctx context.Context) (int64, *helpers.ApiError)

	// Invite user methods
	CreateInviteUser(ctx context.Context, dto *dtos.CreateInviteUserDto) (*models.ProfileResponse, *helpers.ApiError)

	// User profile update
	UpdateMyProfile(ctx context.Context, dto *dtos.CreateProfileDto) (*models.ProfileResponse, *helpers.ApiError)
	GetMe(ctx context.Context) (*models.ProfileResponse, *helpers.ApiError)

	// Migration methods
	MigrateUsersToProfilesDirect(ctx context.Context, limit int64, offset int64) (*types.MigrationResponse, *helpers.ApiError)

	// Common methods
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProfileResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProfilesWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.ProfileResponse, *helpers.ApiError)
	GetProfileByUserId(ctx context.Context, userId primitive.ObjectID) (*models.ProfileResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

// profileService is the implementation of the ProfileService interface
type profileService struct {
	profileRepo ProfileRepo
	usersgw     *gateway.UsersGw
	rolesgw     *gateway.RolesGw
	settingRepo config_services.SettingRepo
	emailgw     *gateway.EmailGw
	db          *mongo.Client
}

// NewProfileService creates a new instance of the ProfileService
func NewProfileService(i *do.Injector) (ProfileService, error) {
	return &profileService{
		profileRepo: do.MustInvoke[ProfileRepo](i),
		usersgw:     do.MustInvoke[*gateway.UsersGw](i),
		rolesgw:     do.MustInvoke[*gateway.RolesGw](i),
		settingRepo: do.MustInvoke[config_services.SettingRepo](i),
		emailgw:     do.MustInvoke[*gateway.EmailGw](i),
		db:          do.MustInvoke[*mongo.Client](i),
	}, nil
}
