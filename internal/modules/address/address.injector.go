package address_init

import (
	address_handlers "larsa-e-commerce-microservice-v2/internal/modules/address/handlers"
	address_repos "larsa-e-commerce-microservice-v2/internal/modules/address/repos"
	address_services "larsa-e-commerce-microservice-v2/internal/modules/address/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

// Init registers all address-related dependencies
func Init(i *do.Injector, r *chi.Mux) *do.Injector {
	// Register the repository
	do.Provide(i, address_repos.NewAddressRepository)

	// Register the service
	do.Provide(i, address_services.NewAddressService)

	// Register the handler and set up routes
	address_handlers.NewAddressHandler(i, r)

	return i
}
