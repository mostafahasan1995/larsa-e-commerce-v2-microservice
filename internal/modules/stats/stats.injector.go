package stats

import (
	stats_handlers "larsa-e-commerce-microservice-v2/internal/modules/stats/handlers"
	stats_repos "larsa-e-commerce-microservice-v2/internal/modules/stats/repos"
	"larsa-e-commerce-microservice-v2/internal/modules/stats/stats_service"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r *chi.Mux) {
	// Repos
	do.Provide(i, stats_repos.NewStatsOrderRepo)
	do.Provide(i, stats_repos.NewStatsInvoiceRepo)

	// Services
	do.Provide(i, stats_service.NewStatsOrderService)

	// Handlers
	stats_handlers.NewStatsHandler(i, r)
}
