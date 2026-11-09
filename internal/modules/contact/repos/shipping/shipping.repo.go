package contact_repos

import (
	"context"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/gateway"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos/shipping"
	filters "larsa-e-commerce-microservice-v2/internal/modules/contact/filters/shipping"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models/shipping"
	services "larsa-e-commerce-microservice-v2/internal/modules/contact/services/shipping"
	utils "larsa-e-commerce-microservice-v2/internal/modules/contact/utils/shipping"
)

type ShippingRepo struct {
	*base_repos.BaseRepo[models.Shipping, models.ShippingResponse]
	usersgw *gateway.UsersGw
	rolesgw *gateway.RolesGw
}

func NewShippingRepo(i *do.Injector) (services.ShippingRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &ShippingRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Shipping, models.ShippingResponse](db, enums.SHIPPING_COMPANY),
		usersgw:  do.MustInvoke[*gateway.UsersGw](i),
		rolesgw:  do.MustInvoke[*gateway.RolesGw](i),
	}, nil
}

func (repo *ShippingRepo) Create(ctx context.Context, createShippingDto *dtos.CreateShippingDto) (*models.ShippingResponse, *helpers.ApiError) {
	//Set an agent role in the user
	role, err := repo.rolesgw.GetRoleByName(ctx, string(enums.ShippingCompanyRoleName))
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "can't find shipping company role",
		}
	}

	//Create a user with the agent role
	//TODO: FIXME, Remove this after testing
	// password := util.GeneratePassword(4, 4, 4, 4)
	password := "123456"
	user, err := repo.usersgw.AddUser(ctx, &types.PostUserData{
		Email:        createShippingDto.Contract.Email,
		FirstName:    createShippingDto.BasicInformation.Name.GetValue(ctx),
		LastName:     createShippingDto.BasicInformation.Name.GetValue(ctx),
		Password:     password,
		Roles:        []primitive.ObjectID{role.Id},
		Capabilities: []primitive.ObjectID{},
	})
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	cfg, _ := util.GetReqAppCfg(ctx)

	shipping := utils.ConvertToShippingModel(createShippingDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	//Set the user id to the shipping
	shipping.Id = user
	createdShipping, createShippingErr := repo.BaseRepo.Create(ctx, *shipping)
	if createShippingErr != nil {
		return nil, createShippingErr
	}

	//TODO: Need an implementation for this from the users service
	// //Send a set password email to the user
	// sendSetPasswordEmailErr := repo.usersgw.SendSetPasswordEmail(ctx, createShippingDto.Contract.Email)
	// if sendSetPasswordEmailErr != nil {
	// 	return nil, &helpers.ApiError{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    sendSetPasswordEmailErr.Error(),
	// 	}
	// }

	return createdShipping, nil
}

func (repo *ShippingRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ShippingResponse, *helpers.ApiError) {
	// Create filters from query parameters
	shippingFilters, _ := filters.NewShippingFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := shippingFilters.BuildPipeline(matchStage, ctx, true)

	shipping, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return shipping, nil
}

func (repo *ShippingRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ShippingsWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	shippingFilters, _ := filters.NewShippingFilters(query)

	// Build pipeline with filters
	pipeline := shippingFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &shippingFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ShippingRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateShippingDto *dtos.CreateShippingDto) (*models.ShippingResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	shipping := utils.ConvertToShippingModel(updateShippingDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *shipping, matchStage)
}

func (repo *ShippingRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}
	return nil
}

func (repo *ShippingRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}

// FindOneByCountryCode finds a shipping company that supports the given ISO3 country code
func (repo *ShippingRepo) FindOneByCountryCode(ctx context.Context, iso3 string) (*models.ShippingResponse, *helpers.ApiError) {
	result, err := repo.BaseRepo.FindOne(ctx, []bson.M{
		{
			"$match": bson.M{
				"contract.countryCodes": bson.M{"$in": []string{iso3}},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
