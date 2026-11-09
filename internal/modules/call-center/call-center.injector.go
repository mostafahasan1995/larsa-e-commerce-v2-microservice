package call_center_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/call-center/handlers"
	services "larsa-e-commerce-microservice-v2/internal/modules/call-center/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// services
	do.Provide(injector, services.NewCallCenterService)

	// handlers
	handlers.NewCallCenterHandler(injector, r)

	return injector
}
