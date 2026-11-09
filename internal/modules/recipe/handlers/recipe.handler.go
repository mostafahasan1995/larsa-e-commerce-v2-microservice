package recipe_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	recipe_dto "larsa-e-commerce-microservice-v2/internal/modules/recipe/dtos"
	recipe_filters "larsa-e-commerce-microservice-v2/internal/modules/recipe/filters"
	recipe_service "larsa-e-commerce-microservice-v2/internal/modules/recipe/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeHandler struct {
	recipeService      recipe_service.RecipeService
	validationInstance *validator.Validate
}

func NewRecipeHandler(i *do.Injector, r *chi.Mux) {
	h := RecipeHandler{
		recipeService:      do.MustInvoke[recipe_service.RecipeService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.RECIPE.GetAdminRoute(), func(r chi.Router) {
		//Create Recipe
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Recipe, enums.Post)),
			middleware.ValidateDTO(recipe_dto.CreateRecipeDTO{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		//Get Recipe By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Recipe, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(recipe_filters.RecipeFilters{}, middleware.QUERY, i),
		).Get("/{id}", h.GetById)

		//Get Recipes With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Recipe, enums.Get)),
			middleware.ValidateDTO(recipe_filters.RecipeFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		//Get All Recipes
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Recipe, enums.Get)),
			middleware.ValidateDTO(recipe_filters.RecipeFilters{}, middleware.QUERY, i),
		).Get("/all", h.GetAll)

		//Update Recipe
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Recipe, enums.Put)),
			middleware.ValidateDTO(recipe_dto.UpdateRecipeDTO{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.Update))

		//Delete Recipe
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Recipe, enums.Delete)),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.Delete)
	})
}

func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, err := middleware.GetValidatedDTO[recipe_dto.CreateRecipeDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}
	result, createErr := h.recipeService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *RecipeHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	query := r.URL.Query().Get("query")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, getErr := h.recipeService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *RecipeHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.recipeService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *RecipeHandler) Update(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	dto, err := middleware.GetValidatedDTO[recipe_dto.UpdateRecipeDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.recipeService.Update(ctx, objectId, dto)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *RecipeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	if deleteErr := h.recipeService.Delete(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RecipeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, getErr := h.recipeService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
