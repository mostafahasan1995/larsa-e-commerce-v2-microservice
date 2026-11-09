package seed_handlers

import (
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	seed_services "larsa-e-commerce-microservice-v2/internal/modules/seed/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type SeedHandler struct {
	seedService seed_services.SeedService
}

func NewSeedHandler(i *do.Injector, r *chi.Mux) {
	h := SeedHandler{
		seedService: do.MustInvoke[seed_services.SeedService](i),
	}

	r.Route(enums.SEED.GetAdminRoute(), func(r chi.Router) {
		// r.With(
		// 	middleware.AdminAuth(middleware.Restriction{
		// 		SuperAdminCapability: enums.CapabilitiesRegistry.General,
		// 		OwnerCapability:      enums.CapabilitiesRegistry.General,
		// 		MemberCapability:     enums.CapabilitiesRegistry.General,
		// 	}),
		// ).Put("/", helpers.Make(h.SeedAll))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
		).Put("/system", helpers.Make(h.SeedSystemData))
		// r.With(
		// 	middleware.AdminAuth(middleware.Restriction{
		// 		SuperAdminCapability: enums.CapabilitiesRegistry.General,
		// 		OwnerCapability:      enums.CapabilitiesRegistry.General,
		// 		MemberCapability:     enums.CapabilitiesRegistry.General,
		// 	}),
		// ).Delete("/", helpers.Make(h.TruncateAllData))
	})
}

func (h *SeedHandler) SeedAll(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	if err := h.seedService.SeedAll(ctx); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *SeedHandler) SeedSystemData(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	if err := h.seedService.SeedSystemData(ctx); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *SeedHandler) TruncateAllData(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	if err := h.seedService.TruncateAllData(ctx); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
