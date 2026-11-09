package contact_repos

import (
	"context"
	"fmt"
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/gateway"
	services "larsa-e-commerce-microservice-v2/internal/modules/contact/services/agent"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos/agent"
	agent_filters "larsa-e-commerce-microservice-v2/internal/modules/contact/filters/agent"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models/agent"
	utils "larsa-e-commerce-microservice-v2/internal/modules/contact/utils/agent"
)

type AgentRepo struct {
	*base_repos.BaseRepo[models.Agent, models.AgentResponse]
	usersgw *gateway.UsersGw
	rolesgw *gateway.RolesGw
}

func NewAgentRepo(i *do.Injector) (services.AgentRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &AgentRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Agent, models.AgentResponse](db, enums.AGENT),
		usersgw:  do.MustInvoke[*gateway.UsersGw](i),
		rolesgw:  do.MustInvoke[*gateway.RolesGw](i),
	}, nil
}

func (repo *AgentRepo) Create(ctx context.Context, createAgentDto *dtos.CreateAgentDto) (*models.AgentResponse, *helpers.ApiError) {
	role, err := repo.rolesgw.GetRoleByName(ctx, string(enums.AgentRoleName))
	if err != nil {
		if err == types.ErrNotFound {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "Role 'agent' not found. Please ensure the role exists in the users service.",
			}
		}
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	//TODO: FIXME, Remove this after testing
	// password := util.GeneratePassword(4, 4, 4, 4)
	password := "123456"
	//Create a user with the agent role
	user, err := repo.usersgw.AddUser(ctx, &types.PostUserData{
		Email:        createAgentDto.BusinessInfo.OfficialEmail,
		FirstName:    createAgentDto.BasicInformation.FirstName,
		LastName:     createAgentDto.BasicInformation.LastName,
		Password:     password,
		Roles:        []primitive.ObjectID{role.Id},
		Capabilities: []primitive.ObjectID{},
	})
	if err != nil {
		fmt.Println("Error adding user: ", err)
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	cfg, _ := util.GetReqAppCfg(ctx)

	agent := utils.ConvertToAgentModel(createAgentDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	//Set the user id to the agent
	agent.Id = user
	createdAgent, createAgentErr := repo.BaseRepo.Create(ctx, *agent)
	if createAgentErr != nil {
		return nil, createAgentErr
	}

	//TODO: Need an implementation for this from the users service
	// //Send a set password email to the user
	// sendSetPasswordEmailErr := repo.usersgw.SendSetPasswordEmail(ctx, createAgentDto.BusinessInfo.OfficialEmail)
	// if sendSetPasswordEmailErr != nil {
	// 	return nil, &helpers.ApiError{
	// 		StatusCode: http.StatusInternalServerError,
	// 		Message:    sendSetPasswordEmailErr.Error(),
	// 	}
	// }

	return createdAgent, nil
}

func (repo *AgentRepo) GetById(ctx context.Context, id primitive.ObjectID, filter *agent_filters.AgentFilters) (*models.AgentResponse, *helpers.ApiError) {
	// Build pipeline with filters
	matchStage := bson.M{}
	pipeline := filter.BuildPipeline(matchStage, ctx, true)

	agent, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return agent, nil
}

func (repo *AgentRepo) GetWithPagination(ctx context.Context, skip, limit int64, filter *agent_filters.AgentFilters) (*models.AgentsWithPagination, *helpers.ApiError) {
	matchStage := bson.M{}

	// Build pipeline with filters
	pipeline := filter.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *AgentRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateAgentDto *dtos.CreateAgentDto) (*models.AgentResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	agent := utils.ConvertToAgentModel(updateAgentDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *agent, matchStage)
}

func (repo *AgentRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {

	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *AgentRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
