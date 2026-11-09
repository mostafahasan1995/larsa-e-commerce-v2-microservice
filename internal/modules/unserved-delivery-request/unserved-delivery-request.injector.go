package unserved_delivery_request_init

import (
	"github.com/go-chi/chi/v5"
	"github.com/samber/do"

	unserved_delivery_request_handlers "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/handlers"
	unserved_delivery_request_repos "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/repos"
	unserved_delivery_request_services "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/services"
)

func Init(injector *do.Injector, r *chi.Mux) {
	do.Provide(injector, unserved_delivery_request_repos.NewUnservedDeliveryRequestRepo)
	do.Provide(injector, unserved_delivery_request_services.NewUnservedDeliveryRequestService)
	unserved_delivery_request_handlers.NewUnservedDeliveryRequestHandler(injector, r)
}
