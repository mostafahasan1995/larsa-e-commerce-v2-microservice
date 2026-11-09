package product_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	product_filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type WebBannerHandler struct {
	webBannerService   services.WebBannerService
	validationInstance *validator.Validate
}

func NewWebBannerHandler(i *do.Injector, r *chi.Mux) {
	h := WebBannerHandler{
		webBannerService:   do.MustInvoke[services.WebBannerService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.BANNER.GetWebRoute(), func(r chi.Router) {

		r.With(
			middleware.ValidateDTO(product_filters.BannerFilters{}, middleware.QUERY, i)).Get("/all", h.GetAll)

	})
}

func (h *WebBannerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.webBannerService.GetAll(ctx, query)
	if err != nil {
		helpers.HandleApiError(w, err)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
