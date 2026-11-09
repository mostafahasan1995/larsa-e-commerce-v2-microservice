package recipe_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	recipe_filters "larsa-e-commerce-microservice-v2/internal/modules/recipe/filters"
	recipe_service "larsa-e-commerce-microservice-v2/internal/modules/recipe/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebRecipeHandler struct {
	webRecipeService   recipe_service.WebRecipeService
	validationInstance *validator.Validate
}

func NewWebRecipeHandler(i *do.Injector, r *chi.Mux) {
	h := WebRecipeHandler{
		webRecipeService:   do.MustInvoke[recipe_service.WebRecipeService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.RECIPE.GetWebRoute(), func(r chi.Router) {
		// Get recipe by ID
		r.With(
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(recipe_filters.RecipeFilters{}, middleware.QUERY, i),
		).Get("/{id}", h.GetById)

		// Get recipes with pagination
		r.With(
			middleware.ValidateDTO(recipe_filters.RecipeFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		// Get all recipes
		r.With(
			middleware.ValidateDTO(recipe_filters.RecipeFilters{}, middleware.QUERY, i),
		).Get("/all", h.GetAll)
	})
}

func (h *WebRecipeHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.webRecipeService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebRecipeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.webRecipeService.GetAll(ctx, query)
	if err != nil {
		helpers.HandleApiError(w, err)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebRecipeHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.webRecipeService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
