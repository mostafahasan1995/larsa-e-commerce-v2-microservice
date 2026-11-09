package address_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	address_dtos "larsa-e-commerce-microservice-v2/internal/modules/address/dtos"
	address_services "larsa-e-commerce-microservice-v2/internal/modules/address/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddressHandler handles HTTP requests for the address module
type AddressHandler struct {
	addressService     address_services.AddressService
	validationInstance *validator.Validate
}

// NewAddressHandler creates a new AddressHandler and sets up the routes
func NewAddressHandler(i *do.Injector, r *chi.Mux) {
	h := AddressHandler{
		addressService:     do.MustInvoke[address_services.AddressService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.ADDRESS.GetWebRoute(), func(r chi.Router) {
		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(address_dtos.AddressDTO{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		r.With(
			middleware.Auth("authenticate"),
		).Get("/all", h.GetAll)

		r.With(
			middleware.Auth("authenticate"),
		).Get("/", h.GetWithPagination)

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(address_dtos.AddressDTO{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.DeleteById)

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.BODY, i),
		).Delete("/", helpers.Make(h.BulkDelete))

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Post("/{id}/default", h.SetDefault)
	})
}

// Create handles the creation of a new address
func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[address_dtos.AddressDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.addressService.Create(ctx, *dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

// GetAll handles the retrieval of all addresses
func (h *AddressHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	results, getErr := h.addressService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, results)
	w.Write(response)
}

// GetWithPagination handles the retrieval of addresses with pagination
func (h *AddressHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, paginateErr := util.Paginate(r)
	if paginateErr != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    paginateErr.Error(),
		})
		return
	}

	results, getErr := h.addressService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, results)
	w.Write(response)
}

// GetById handles the retrieval of an address by ID
func (h *AddressHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.addressService.GetById(ctx, objectId)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// UpdateById handles the update of an address by ID
func (h *AddressHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[address_dtos.AddressDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.addressService.UpdateById(ctx, objectId, *data)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

// DeleteById handles the deletion of an address by ID
func (h *AddressHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.addressService.DeleteById(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// BulkDelete handles the deletion of multiple addresses
func (h *AddressHandler) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[types.BulkDeleteDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if bulkDeleteErr := h.addressService.BulkDelete(ctx, data.Ids); bulkDeleteErr != nil {
		helpers.HandleApiError(w, bulkDeleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// SetDefault sets an address as the default for the current user
func (h *AddressHandler) SetDefault(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if err := h.addressService.SetDefaultAddress(ctx, objectId); err != nil {
		helpers.HandleApiError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
