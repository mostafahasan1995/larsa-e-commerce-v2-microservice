package survey_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	survey_filters "larsa-e-commerce-microservice-v2/internal/modules/survey/filters"
	survey_services "larsa-e-commerce-microservice-v2/internal/modules/survey/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebQuestionHandler struct {
	questionService    survey_services.WebQuestionService
	validationInstance *validator.Validate
}

func NewWebQuestionHandler(i *do.Injector, r *chi.Mux) {
	h := WebQuestionHandler{
		questionService:    do.MustInvoke[survey_services.WebQuestionService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.QUESTION.GetWebRoute(), func(r chi.Router) {
		// Get question by ID
		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		// Get questions with pagination
		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(survey_filters.QuestionFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		// Get all questions
		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(survey_filters.QuestionFilters{}, middleware.QUERY, i),
		).Get("/all", h.GetAll)
	})
}

func (h *WebQuestionHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.questionService.GetById(ctx, objectId)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebQuestionHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.questionService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebQuestionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	result, getErr := h.questionService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
