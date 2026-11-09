package settings_inventory

import (
	settings_inventory_handlers "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/handlers"
	settings_inventory_repos "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/repos"
	settings_inventory_services "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r chi.Router) {
	do.Provide(i, settings_inventory_repos.NewInventorySettingsRepository)

	do.Provide(i, settings_inventory_services.NewInventorySettingsService)

	settings_inventory_handlers.NewInventorySettingsHandler(i, r)
}
