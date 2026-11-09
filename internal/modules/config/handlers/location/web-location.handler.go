package config_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_location_services "larsa-e-commerce-microservice-v2/internal/modules/config/services/location"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type LocationWebHandler struct {
	locationWebService config_location_services.LocationWebService
}

func NewLocationWebHandler(i *do.Injector, r *chi.Mux) {
	h := LocationWebHandler{
		locationWebService: do.MustInvoke[config_location_services.LocationWebService](i),
	}

	r.Route("/", func(r chi.Router) {

		//Get All Countries
		r.With().Get("/countries/all", helpers.Make(h.GetAllCountries))

		//Get All Cities
		r.With().Get("/cities/all", helpers.Make(h.GetAllCities))

		//Get All Areas
		r.With().Get("/areas/all", helpers.Make(h.GetAllAreas))
	})
}

func (h *LocationWebHandler) GetAllCountries(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.locationWebService.GetAllCountries(ctx)
	if err != nil {
		return err
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *LocationWebHandler) GetAllCities(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")
	result, err := h.locationWebService.GetCitiesByCountry(ctx, query)
	if err != nil {
		return err
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *LocationWebHandler) GetAllAreas(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")
	result, err := h.locationWebService.GetAreasByCity(ctx, query)
	if err != nil {
		return err
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}
