package loyalty_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/loyalty/handlers"
	repos "larsa-e-commerce-microservice-v2/internal/modules/loyalty/repos"
	services "larsa-e-commerce-microservice-v2/internal/modules/loyalty/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewWebLoyaltyPointsRepo)
	do.Provide(injector, repos.NewLoyaltyOrderRepo)

	// services
	do.Provide(injector, services.NewWebLoyaltyPointsService)
	do.Provide(injector, services.NewLoyaltyOrderService)
	do.Provide(injector, services.NewLoyaltyDiscountService)

	// handlers
	handlers.NewWebLoyaltyPointsHandler(injector, r)
	handlers.NewLoyaltyDiscountHandler(injector, r)
	handlers.NewLoyaltyOrderHandler(injector, r)

	return injector
}
