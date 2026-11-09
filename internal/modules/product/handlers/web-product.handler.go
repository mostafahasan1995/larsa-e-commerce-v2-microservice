package product_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type WebProductHandler struct {
	webProductService  services.WebProductService
	webReviewService   services.WebReviewService
	branchService      services.BranchService
	validationInstance *validator.Validate
}

func NewWebProductHandler(i *do.Injector, r *chi.Mux) {
	h := WebProductHandler{
		webProductService:  do.MustInvoke[services.WebProductService](i),
		webReviewService:   do.MustInvoke[services.WebReviewService](i),
		branchService:      do.MustInvoke[services.BranchService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.PRODUCT.GetWebRoute(), func(r chi.Router) {

		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(product_filters.ProductFilters{}, middleware.QUERY, i)).Get("/{id}", h.GetById)
		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(product_filters.ProductFilters{}, middleware.QUERY, i)).Get("/", h.GetWithPagination)

		r.Route("/{id}"+enums.REVIEW.GetWebRoute(), func(r chi.Router) {
			r.With(
				middleware.Auth([]string{}...),
				middleware.ValidateDTO(dtos.CreateReviewDto{}, middleware.BODY, i),
			).Post("/", helpers.Make(h.CreateReview))
			r.With(
				middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Product, enums.Get)),
				middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
				middleware.ValidateDTO(product_filters.WebReviewFilters{}, middleware.QUERY, i)).Get("/", h.GetReviewsWithPagination)

			r.With(
				middleware.Auth([]string{}...),
				middleware.ValidateDTO(dtos.CreateReviewDto{}, middleware.BODY, i),
				middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			).Put("/", helpers.Make(h.UpdateReviewByProductId))
			r.With(
				middleware.Auth([]string{}...),
				middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			).Delete("/", h.DeleteReviewByProductId)

		})
	})
}

func (h *WebProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	branchId := r.URL.Query().Get("branchId")
	var unitId = r.URL.Query().Get("unitId")
	currencyCode := r.URL.Query().Get("currencyCode")
	if currencyCode == "" {
		currencyCode = configuration.EnvironmentVariables.DefaultCurrencyCode
	}

	var branchObjectId *primitive.ObjectID
	var unitObjectId *primitive.ObjectID

	if branchId != "" {
		objectId, err := primitive.ObjectIDFromHex(branchId)
		if err != nil {
			helpers.DefaultError(err)
			return
		}
		branchObjectId = &objectId
	}

	if unitId != "" {
		objectId, err := primitive.ObjectIDFromHex(unitId)
		if err != nil {
			helpers.DefaultError(err)
			return
		}
		unitObjectId = &objectId
	}

	if r.URL.Query().Get("currencyCode") != "" {
		currencyCode = r.URL.Query().Get("currencyCode")
	}

	result, getErr := h.webProductService.GetById(ctx, product_utils.GetByIdArgs{
		Id:           objectId,
		Query:        query,
		BranchId:     branchObjectId,
		UnitId:       unitObjectId,
		CurrencyCode: &currencyCode,
	})
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	result = GetPrice(ctx, GetPriceArgs{
		Product:       *result,
		BranchId:      branchId,
		UnitId:        unitId,
		CurrencyCode:  currencyCode,
		CalculateList: r.URL.Query().Get("calculatePriceList") == "true",
	})
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebProductHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.HandleApiError(w, helpers.DefaultError(errGetPaginate))
		return
	}

	var branchId = r.URL.Query().Get("branchId")
	var unitId = r.URL.Query().Get("unitId")
	currencyCode := r.URL.Query().Get("currencyCode")
	if currencyCode == "" {
		currencyCode = configuration.EnvironmentVariables.DefaultCurrencyCode
	}

	var branchObjectId *primitive.ObjectID
	var unitObjectId *primitive.ObjectID

	if r.URL.Query().Get("currencyCode") != "" {
		currencyCode = r.URL.Query().Get("currencyCode")
	}

	if branchId != "" {
		objectId, err := primitive.ObjectIDFromHex(branchId)
		if err != nil {
			apiErr := helpers.BadRequest("invalid branch id")
			helpers.HandleApiError(w, &apiErr)
			return
		}
		branchObjectId = &objectId
	}

	if unitId != "" {
		objectId, err := primitive.ObjectIDFromHex(unitId)
		if err != nil {
			apiErr := helpers.BadRequest("invalid unit id")
			helpers.HandleApiError(w, &apiErr)
			return
		}
		unitObjectId = &objectId
	}

	result, getErr := h.webProductService.GetWithPagination(ctx, product_utils.GetWithPaginationArgs{
		Skip:         skip,
		Limit:        limit,
		Query:        query,
		BranchId:     branchObjectId,
		UnitId:       unitObjectId,
		CurrencyCode: &currencyCode,
	})
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	if len(result.Data) == 0 {
		log.Info().
			Str("query", query).
			Str("user_agent", r.UserAgent()).
			Msg("No products found")
	}

	for i, product := range result.Data {
		result.Data[i] = *GetPrice(ctx, GetPriceArgs{
			Product:       product,
			BranchId:      r.URL.Query().Get("branchId"),
			UnitId:        r.URL.Query().Get("unitId"),
			CurrencyCode:  r.URL.Query().Get("currencyCode"),
			CalculateList: r.URL.Query().Get("calculatePriceList") == "true",
		})
		//Remove prices to reduce the response size
		result.Data[i].Prices = []product_models.Price{}
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebProductHandler) CreateReview(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateReviewDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, createErr := h.webReviewService.Create(ctx, objectId, *dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *WebProductHandler) GetReviewsWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.webReviewService.GetByProductIdWithPagination(ctx, objectId, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebProductHandler) UpdateReviewByProductId(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.CreateReviewDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.webReviewService.UpdateByProductId(ctx, objectId, *data)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *WebProductHandler) DeleteReviewByProductId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.webReviewService.DeleteByProductId(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
