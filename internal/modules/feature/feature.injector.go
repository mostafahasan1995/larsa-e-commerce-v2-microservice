package feature_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/feature/handler"
	repos "larsa-e-commerce-microservice-v2/internal/modules/feature/repo"
	services "larsa-e-commerce-microservice-v2/internal/modules/feature/service"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewFeatureRepository)
	do.Provide(injector, repos.NewFeatureGroupRepository)

	// services
	do.Provide(injector, services.NewFeatureService)
	do.Provide(injector, services.NewFeatureGroupService)

	// handlers
	handlers.NewFeatureHandler(injector, r)
	handlers.NewFeatureGroupHandler(injector, r)

	return injector
}
