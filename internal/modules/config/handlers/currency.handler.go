package config_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CurrencyHandler struct {
	currencyService    services.CurrencyService
	validationInstance *validator.Validate
}

func NewCurrencyHandler(i *do.Injector, r *chi.Mux) {
	h := CurrencyHandler{
		currencyService:    do.MustInvoke[services.CurrencyService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.CURRENCY.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Currency.SuperAdmin.Post.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Currency.Vendor.Post.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Currency.Member.Post.String(),
			}),
			middleware.ValidateDTO(dtos.CreateCurrencyDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Currency.SuperAdmin.GetOne.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Currency.Vendor.GetOne.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Currency.Member.GetOne.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(config_filters.CurrencyFilters{}, middleware.QUERY, i)).Get("/{id}", h.GetById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(config_filters.CurrencyFilters{}, middleware.QUERY, i)).Get("/all", h.GetAll)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(config_filters.CurrencyFilters{}, middleware.QUERY, i)).Get("/", h.GetWithPagination)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Currency.SuperAdmin.Put.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Currency.Vendor.Put.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Currency.Member.Put.String(),
			}),
			middleware.ValidateDTO(dtos.CreateCurrencyDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))

		// r.With(
		// 	middleware.AdminAuth(middleware.Restriction{
		// 		SuperAdminCapability: enums.CapabilitiesRegistry.Currency.SuperAdmin.Delete.String(),
		// 		OwnerCapability:      enums.CapabilitiesRegistry.Currency.Vendor.Delete.String(),
		// 		MemberCapability:     enums.CapabilitiesRegistry.Currency.Member.Delete.String(),
		// 	}),
		// 	middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		// ).Delete("/{id}", h.DeleteById)
		// r.With(
		// 	middleware.AdminAuth(middleware.Restriction{
		// 		SuperAdminCapability: enums.CapabilitiesRegistry.Currency.SuperAdmin.DeleteMany.String(),
		// 		OwnerCapability:      enums.CapabilitiesRegistry.Currency.Vendor.DeleteMany.String(),
		// 		MemberCapability:     enums.CapabilitiesRegistry.Currency.Member.DeleteMany.String(),
		// 	}),
		// 	middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.BODY, i),
		// ).Delete("/", helpers.Make(h.BulkDelete))
	})
}

func (h *CurrencyHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateCurrencyDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.currencyService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *CurrencyHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.currencyService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *CurrencyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.currencyService.GetAll(ctx, query)
	if err != nil {
		helpers.HandleApiError(w, err)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *CurrencyHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.currencyService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *CurrencyHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.CreateCurrencyDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.currencyService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *CurrencyHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.currencyService.DeleteById(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CurrencyHandler) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[types.BulkDeleteDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if bulkDeleteErr := h.currencyService.BulkDelete(ctx, data.Ids); bulkDeleteErr != nil {
		helpers.HandleApiError(w, bulkDeleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
