package config_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type MigrationHandler struct {
	migrationService services.MigrationService
}

func NewMigrationHandler(i *do.Injector, r *chi.Mux) {
	h := MigrationHandler{
		migrationService: do.MustInvoke[services.MigrationService](i),
	}

	r.Route("/admin/migration", func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Category.SuperAdmin.Put.String(),
			}),
			middleware.ValidateDTO(config_dtos.DataMigrationDto{}, middleware.BODY, i),
		).Post("/migrate-data", helpers.Make(h.MigrateData))
	})
}

func (h *MigrationHandler) MigrateData(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[config_dtos.DataMigrationDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if migrateErr := h.migrationService.MigrateData(ctx, *dto); migrateErr != nil {
		helpers.HandleApiError(w, migrateErr)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	helpers.WriteJson(w, http.StatusOK, map[string]string{
		"message": "Data migration completed successfully",
	})

	return nil
}
