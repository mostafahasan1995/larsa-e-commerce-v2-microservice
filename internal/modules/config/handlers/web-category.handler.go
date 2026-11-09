package config_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebCategoryHandler struct {
	webCategoryService services.WebCategoryService
	validationInstance *validator.Validate
}

func NewWebCategoryHandler(i *do.Injector, r *chi.Mux) {
	h := WebCategoryHandler{
		webCategoryService: do.MustInvoke[services.WebCategoryService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.CATEGORY.GetWebRoute(), func(r chi.Router) {
		r.With(

			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(config_filters.CategoryFilters{}, middleware.QUERY, i)).Get("/{id}", h.GetById)
		r.With(

			middleware.ValidateDTO(config_filters.CategoryFilters{}, middleware.QUERY, i)).Get("/all", h.GetAll)
		r.With(

			middleware.ValidateDTO(config_filters.CategoryFilters{}, middleware.QUERY, i)).Get("/", h.GetWithPagination)

	})
}

func (h *WebCategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.webCategoryService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.webCategoryService.GetAll(ctx, query)
	if err != nil {
		helpers.HandleApiError(w, err)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebCategoryHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.webCategoryService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
