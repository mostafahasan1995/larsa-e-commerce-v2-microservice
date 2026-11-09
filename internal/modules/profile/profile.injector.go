package profile_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/profile/handlers"
	repos "larsa-e-commerce-microservice-v2/internal/modules/profile/repos"
	services "larsa-e-commerce-microservice-v2/internal/modules/profile/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

// Init initializes the profile module
func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// Register repos
	do.Provide(injector, repos.NewProfileRepo)

	// Register services
	do.Provide(injector, services.NewProfileService)

	// Register handlers
	handlers.NewProfileHandler(injector, r)

	return injector
}
