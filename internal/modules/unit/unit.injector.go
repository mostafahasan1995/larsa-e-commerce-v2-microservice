package unit_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/unit/handlers"
	repos "larsa-e-commerce-microservice-v2/internal/modules/unit/repos"
	services "larsa-e-commerce-microservice-v2/internal/modules/unit/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewUnitRepo)
	do.Provide(injector, repos.NewWebUnitRepo)

	// services
	do.Provide(injector, services.NewUnitService)
	do.Provide(injector, services.NewWebUnitService)

	// handlers
	handlers.NewUnitHandler(injector, r)
	handlers.NewWebUnitHandler(injector, r)

	return injector
}
