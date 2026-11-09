package settings_sales

import (
	settings_sales_handlers "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/handlers"
	settings_sales_repos "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/repos"
	settings_sales_services "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r chi.Router) {
	do.Provide(i, settings_sales_repos.NewSalesSettingsRepository)
	do.Provide(i, settings_sales_repos.NewSalesSettingsWebRepository)

	do.Provide(i, settings_sales_services.NewSalesSettingsService)
	do.Provide(i, settings_sales_services.NewSalesSettingsWebService)

	settings_sales_handlers.NewSalesSettingsHandler(i, r)
	settings_sales_handlers.NewSalesSettingsWebHandler(i, r)
}
