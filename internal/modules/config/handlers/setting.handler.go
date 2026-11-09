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

type SettingHandler struct {
	settingService services.SettingService
}

func NewSettingHandler(i *do.Injector, r *chi.Mux) {
	h := SettingHandler{
		settingService: do.MustInvoke[services.SettingService](i),
	}

	r.Route(enums.SETTING.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Category.SuperAdmin.Put.String(),
			}),
			middleware.ValidateDTO(config_dtos.TruncateCollectionDtoList{}, middleware.BODY, i),
		).Delete("/truncate-collections", helpers.Make(h.TruncateCollections))

	})
}

func (h *SettingHandler) TruncateCollections(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	cfg, _ := util.GetReqAppCfg(ctx)
	dto, err := middleware.GetValidatedDTO[config_dtos.TruncateCollectionDtoList](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if deleteErr := h.settingService.TruncateCollections(cfg.Db, *dto); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
