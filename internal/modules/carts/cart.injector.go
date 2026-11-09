package cart_init

import (
	"github.com/go-chi/chi/v5"
	"github.com/samber/do"

	cart_handlers "larsa-e-commerce-microservice-v2/internal/modules/carts/handlers"
	cart_repos "larsa-e-commerce-microservice-v2/internal/modules/carts/repos"
	cart_services "larsa-e-commerce-microservice-v2/internal/modules/carts/services"
)

// ProvideServices registers all cart services and repos without creating handlers
func ProvideServices(injector *do.Injector) {
	// repos
	do.Provide(injector, cart_repos.NewCartItemRepo)
	do.Provide(injector, cart_repos.NewCartRepo)
	do.Provide(injector, cart_repos.NewCheckoutRepo)

	// services
	do.Provide(injector, cart_services.NewCartItemService)
	do.Provide(injector, cart_services.NewCartService)
	do.Provide(injector, cart_services.NewCheckoutService)
}

// InitHandlers creates handlers after all services are registered
func InitHandlers(injector *do.Injector, r *chi.Mux) {
	// handlers
	cart_handlers.NewCartItemHandler(injector, r)
	cart_handlers.NewCheckoutHandler(injector, r)
}

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	ProvideServices(injector)
	InitHandlers(injector, r)
	return injector
}
