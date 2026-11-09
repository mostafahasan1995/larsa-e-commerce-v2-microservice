package settings_sales_handlers

import (
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	settings_sales_services "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type SalesSettingsWebHandler struct {
	service settings_sales_services.SalesSettingsWebService
}

func NewSalesSettingsWebHandler(i *do.Injector, r chi.Router) {
	h := &SalesSettingsWebHandler{
		service: do.MustInvoke[settings_sales_services.SalesSettingsWebService](i),
	}

	r.Route(enums.SALES_SETTINGS.GetWebRoute(), func(r chi.Router) {
		r.Get("/", helpers.Make(h.GetSalesSettings))
	})
}

// notes
func (h *SalesSettingsWebHandler) GetSalesSettings(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	settings, apiErr := h.service.GetSalesSettings(ctx)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, settings)
}
