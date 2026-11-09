package ticket

import (
	ticket_handlers "larsa-e-commerce-microservice-v2/internal/modules/ticket/handlers"
	ticket_repos "larsa-e-commerce-microservice-v2/internal/modules/ticket/repos"
	ticket_service "larsa-e-commerce-microservice-v2/internal/modules/ticket/service"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) {
	do.Provide(injector, ticket_repos.NewTicketRepository)

	do.Provide(injector, ticket_service.NewTicketService)

	ticket_handlers.NewTicketHandler(injector, r)
}
