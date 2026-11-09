package recipe_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/recipe/handlers"
	repos "larsa-e-commerce-microservice-v2/internal/modules/recipe/repos"
	services "larsa-e-commerce-microservice-v2/internal/modules/recipe/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewRecipeRepository)
	do.Provide(injector, repos.NewWebRecipeRepository)

	// services
	do.Provide(injector, services.NewRecipeService)
	do.Provide(injector, services.NewWebRecipeService)

	// handlers
	handlers.NewRecipeHandler(injector, r)
	handlers.NewWebRecipeHandler(injector, r)

	return injector
}
