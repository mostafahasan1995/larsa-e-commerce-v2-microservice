package quotation_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_utils "larsa-e-commerce-microservice-v2/internal/modules/quotation/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type QuotationAdminSettingsHandler struct {
}

func NewQuotationAdminSettingsHandler(i *do.Injector, r *chi.Mux) {
	h := QuotationAdminSettingsHandler{}

	r.Route(enums.QUOTATION_SETTINGS.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.QuotationSettings, enums.Get)),
			middleware.ValidateDTO(quotation_dto.QuotationSettingsDTO{}, middleware.QUERY, i),
		).Get("/", helpers.Make(h.GetQuotationSettings))
	})
}

func (h *QuotationAdminSettingsHandler) GetQuotationSettings(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Get default quotation settings
	quotationSettings := quotation_utils.GetDefaultQuotationSettings()

	response, _ := json.MarshalContext(ctx, quotationSettings)
	w.Write(response)
	return nil
}
