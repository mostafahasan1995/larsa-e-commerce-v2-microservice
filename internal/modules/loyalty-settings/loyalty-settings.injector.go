package loyalty_settings

import (
	loyalty_settings_handlers "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/handlers"
	loyalty_settings_repos "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/repos"
	loyalty_settings_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func RegisterLoyaltySettingsModule(i *do.Injector, r *chi.Mux) {
	// Register repositories
	do.Provide(i, loyalty_settings_repos.NewLoyaltySettingsRepo)

	// Register services
	do.Provide(i, loyalty_settings_services.NewLoyaltySettingsService)

	// Register handlers
	loyalty_settings_handlers.NewLoyaltySettingsHandler(i, r)
}
