package quotation_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_filters "larsa-e-commerce-microservice-v2/internal/modules/quotation/filters"
	quotation_service "larsa-e-commerce-microservice-v2/internal/modules/quotation/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type QuotationBuyerHandler struct {
	quotationBuyerService quotation_service.QuotationBuyerService
	adminService          quotation_service.QuotationAdminService
	validationInstance    *validator.Validate
}

func NewQuotationBuyerHandler(i *do.Injector, r *chi.Mux) {
	h := QuotationBuyerHandler{
		quotationBuyerService: do.MustInvoke[quotation_service.QuotationBuyerService](i),
		adminService:          do.MustInvoke[quotation_service.QuotationAdminService](i),
		validationInstance:    do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.QUOTATION.GetWebRoute(), func(r chi.Router) {

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(quotation_dto.CreateQuotationDTO{}, middleware.BODY, i),
		).Post("/create-quotation-request", helpers.Make(h.CreateRequestQuotation))

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(quotation_filters.QuotationFilters{}, middleware.QUERY, i),
		).Get("/", helpers.Make(h.GetQuotations))

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(quotation_filters.QuotationFilters{}, middleware.QUERY, i),
		).Get("/general", helpers.Make(h.GetGeneralQuotations))

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(quotation_dto.QuotationBuyerInfoDTO{}, middleware.BODY, i),
		).Post("/general/{id}", helpers.Make(h.RequestGeneralQuotation))

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(quotation_filters.QuotationFilters{}, middleware.QUERY, i),
		).Get("/{id}/response", helpers.Make(h.GetResponseQuotations))

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", helpers.Make(h.GetQuotationById))

		r.With(
			middleware.Auth("authenticate"),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Post("/{id}/status", helpers.Make(h.UpdateQuotationStatus))

	})
}

func (h *QuotationBuyerHandler) CreateRequestQuotation(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, _ := middleware.GetValidatedDTO[quotation_dto.CreateQuotationDTO](r)

	result, apiErr := h.quotationBuyerService.CreateRequestQuotation(ctx, dto)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	return nil
}

func (h *QuotationBuyerHandler) GetQuotationById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	query := r.URL.Query().Get("query")

	result, getErr := h.quotationBuyerService.GetById(ctx, id, query)
	if getErr != nil {
		return getErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

func (h *QuotationBuyerHandler) GetResponseQuotations(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	result, getErr := h.quotationBuyerService.GetResponseQuotations(ctx, skip, limit, id)
	if getErr != nil {
		return getErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

func (h *QuotationBuyerHandler) GetGeneralQuotations(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	result, getErr := h.quotationBuyerService.GetGeneralQuotations(ctx, skip, limit, query)
	if getErr != nil {
		return getErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

func (h *QuotationBuyerHandler) GetQuotations(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	result, getErr := h.quotationBuyerService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		return getErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

func (h *QuotationBuyerHandler) UpdateQuotationStatus(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	status := quotation_enums.QuotationStatus(r.URL.Query().Get("status"))

	apiErr := h.adminService.UpdateQuotationStatus(ctx, id, status)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, nil)
	w.Write(response)
	return nil
}

func (h *QuotationBuyerHandler) RequestGeneralQuotation(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	dto, _ := middleware.GetValidatedDTO[quotation_dto.QuotationBuyerInfoDTO](r)

	result, apiErr := h.quotationBuyerService.RequestGeneralQuotation(ctx, id, *dto)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}
