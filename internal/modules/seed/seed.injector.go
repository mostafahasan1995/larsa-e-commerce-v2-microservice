package seed_init

import (
	"github.com/go-chi/chi/v5"
	"github.com/samber/do"

	seed_handlers "larsa-e-commerce-microservice-v2/internal/modules/seed/handlers"
	seed_services "larsa-e-commerce-microservice-v2/internal/modules/seed/services"
)

func Init(i *do.Injector, r *chi.Mux) {
	do.Provide(i, seed_services.NewSeedService)
	seed_handlers.NewSeedHandler(i, r)
}
