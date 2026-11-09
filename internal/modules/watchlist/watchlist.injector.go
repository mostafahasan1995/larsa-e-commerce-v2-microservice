package watchlist_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/watchlist/handlers"
	repos "larsa-e-commerce-microservice-v2/internal/modules/watchlist/repos"
	services "larsa-e-commerce-microservice-v2/internal/modules/watchlist/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	do.Provide(injector, repos.NewWebWatchListRepo)

	do.Provide(injector, services.NewWebWatchListService)

	handlers.NewWebWatchListHandler(injector, r)

	return injector
}
