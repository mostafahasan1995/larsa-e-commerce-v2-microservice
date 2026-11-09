package plan_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/plan/handler"
	repos "larsa-e-commerce-microservice-v2/internal/modules/plan/repos"
	services "larsa-e-commerce-microservice-v2/internal/modules/plan/service"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewPlanRepository)
	do.Provide(injector, repos.NewWebPlanRepo)

	// services
	do.Provide(injector, services.NewPlanService)
	do.Provide(injector, services.NewWebPlanService)

	// handlers
	handlers.NewPlanHandler(injector, r)
	handlers.NewWebPlanHandler(injector, r)

	return injector
}
