package config_init

import (
	handlers "larsa-e-commerce-microservice-v2/internal/modules/config/handlers"
	location_handlers "larsa-e-commerce-microservice-v2/internal/modules/config/handlers/location"
	repos "larsa-e-commerce-microservice-v2/internal/modules/config/repos"
	location_repos "larsa-e-commerce-microservice-v2/internal/modules/config/repos/location"
	services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	location_services "larsa-e-commerce-microservice-v2/internal/modules/config/services/location"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, repos.NewCurrencyRepo)
	do.Provide(injector, repos.NewCategoryRepo)
	do.Provide(injector, repos.NewWebCategoryRepo)
	do.Provide(injector, repos.NewSettingRepo)
	do.Provide(injector, repos.NewZoneRepo)
	do.Provide(injector, repos.NewTicketTypeRepo)
	do.Provide(injector, repos.NewMigrationRepo)

	//location repos
	do.Provide(injector, location_repos.NewCountryRepo)
	do.Provide(injector, location_repos.NewCityRepo)
	do.Provide(injector, location_repos.NewAreaRepo)

	// services
	do.Provide(injector, services.NewCurrencyService)
	do.Provide(injector, services.NewCategoryService)
	do.Provide(injector, services.NewWebCategoryService)
	do.Provide(injector, services.NewCapabilitiesService)
	do.Provide(injector, services.NewSettingService)
	do.Provide(injector, services.NewZoneService)
	do.Provide(injector, services.NewTicketTypeService)
	do.Provide(injector, services.NewMigrationService)

	//location services
	do.Provide(injector, location_services.NewCountryService)
	do.Provide(injector, location_services.NewCityService)
	do.Provide(injector, location_services.NewAreaService)
	do.Provide(injector, location_services.NewLocationWebService)

	// handlers
	handlers.NewCurrencyHandler(injector, r)
	handlers.NewCategoryHandler(injector, r)
	handlers.NewWebCategoryHandler(injector, r)
	handlers.NewCapabilitiesHandler(injector, r)
	handlers.NewZoneHandler(injector, r)
	handlers.NewTicketTypeHandler(injector, r)
	handlers.NewSettingHandler(injector, r)
	handlers.NewMigrationHandler(injector, r)

	//location handlers
	location_handlers.NewCountryHandler(injector, r)
	location_handlers.NewCityHandler(injector, r)
	location_handlers.NewAreaHandler(injector, r)
	location_handlers.NewLocationWebHandler(injector, r)

	return injector
}
