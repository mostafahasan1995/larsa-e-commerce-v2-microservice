package contact_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/contact/handlers"
	agent_handlers "larsa-e-commerce-microservice-v2/internal/modules/contact/handlers/agent"
	shipping_handlers "larsa-e-commerce-microservice-v2/internal/modules/contact/handlers/shipping"
	repos "larsa-e-commerce-microservice-v2/internal/modules/contact/repos"
	agent_repos "larsa-e-commerce-microservice-v2/internal/modules/contact/repos/agent"
	shipping_repos "larsa-e-commerce-microservice-v2/internal/modules/contact/repos/shipping"
	services "larsa-e-commerce-microservice-v2/internal/modules/contact/services"
	agent_services "larsa-e-commerce-microservice-v2/internal/modules/contact/services/agent"
	shipping_services "larsa-e-commerce-microservice-v2/internal/modules/contact/services/shipping"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewBrokerRepo)
	do.Provide(injector, agent_repos.NewAgentRepo)
	do.Provide(injector, shipping_repos.NewShippingRepo)
	do.Provide(injector, repos.NewCustomerRepo)

	// services
	do.Provide(injector, services.NewBrokerService)
	do.Provide(injector, agent_services.NewAgentService)
	do.Provide(injector, shipping_services.NewShippingService)
	do.Provide(injector, services.NewCustomerService)

	// handlers
	handlers.NewBrokerHandler(injector, r)
	agent_handlers.NewAgentHandler(injector, r)
	shipping_handlers.NewShippingHandler(injector, r)
	handlers.NewCustomerHandler(injector, r)

	return injector
}
