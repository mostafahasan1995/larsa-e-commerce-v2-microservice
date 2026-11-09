package profile_repos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/gateway"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/profile/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProfileRepo is the implementation of the ProfileRepo interface
type ProfileRepo struct {
	*base_repos.BaseRepo[models.CommonProfile, models.ProfileResponse]
	usersgw *gateway.UsersGw
}

// NewProfileRepo creates a new instance of ProfileRepo
func NewProfileRepo(i *do.Injector) (services.ProfileRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &ProfileRepo{
		BaseRepo: base_repos.NewBaseRepo[models.CommonProfile, models.ProfileResponse](db, enums.PROFILE),
		usersgw:  do.MustInvoke[*gateway.UsersGw](i),
	}, nil
}

// NOTE: Methods are split across multiple files in this package for maintainability.
