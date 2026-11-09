package appointment_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/appointment/handler"
	repos "larsa-e-commerce-microservice-v2/internal/modules/appointment/repos"
	services "larsa-e-commerce-microservice-v2/internal/modules/appointment/service"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewAppointmentRepository)

	// services
	do.Provide(injector, services.NewAppointmentService)

	// handlers
	handlers.NewAppointmentHandler(injector, r)

	return injector
}
