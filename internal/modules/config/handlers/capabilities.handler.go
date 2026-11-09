package config_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	"net/http"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type CapabilitiesHandler struct {
	capabilitiesService config_services.CapabilitiesService
}

func NewCapabilitiesHandler(i *do.Injector, r *chi.Mux) {
	h := CapabilitiesHandler{
		capabilitiesService: do.MustInvoke[config_services.CapabilitiesService](i),
	}

	r.Route(enums.CAPABILITY.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: "authenticate",
			}),
		).Post("/setup", helpers.Make(h.SetupCapabilities))
	})
}

func (h *CapabilitiesHandler) SetupCapabilities(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	headerParams := common.ExtractHeaderParams(r)

	if err := h.capabilitiesService.SetupCapabilities(ctx, headerParams); err != nil {
		helpers.HandleApiError(w, helpers.DefaultError(err))
		return nil
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Capabilities setup completed"}`))
	return nil
}
