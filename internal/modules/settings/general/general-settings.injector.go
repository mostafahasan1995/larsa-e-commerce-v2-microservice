package settings_general

import (
	settings_general_handlers "larsa-e-commerce-microservice-v2/internal/modules/settings/general/handlers"
	settings_general_repos "larsa-e-commerce-microservice-v2/internal/modules/settings/general/repos"
	settings_general_services "larsa-e-commerce-microservice-v2/internal/modules/settings/general/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(i *do.Injector, r chi.Router) {
	do.Provide(i, settings_general_repos.NewGeneralSettingsRepository)

	do.Provide(i, settings_general_services.NewGeneralSettingsService)

	settings_general_handlers.NewGeneralSettingsHandler(i, r)
}
