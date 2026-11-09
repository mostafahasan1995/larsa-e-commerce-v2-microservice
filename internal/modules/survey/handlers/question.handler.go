package survey_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	survey_dtos "larsa-e-commerce-microservice-v2/internal/modules/survey/dtos"
	survey_filters "larsa-e-commerce-microservice-v2/internal/modules/survey/filters"
	survey_services "larsa-e-commerce-microservice-v2/internal/modules/survey/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionHandler struct {
	questionService    survey_services.QuestionService
	validationInstance *validator.Validate
}

func NewQuestionHandler(i *do.Injector, r *chi.Mux) {
	h := QuestionHandler{
		questionService:    do.MustInvoke[survey_services.QuestionService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.QUESTION.GetAdminRoute(), func(r chi.Router) {
		//Create Question
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Question, enums.Post)),
			middleware.ValidateDTO(survey_dtos.CreateQuestionDTO{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		//Get Question By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Question, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		//Get Questions With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Question, enums.Get)),
			middleware.ValidateDTO(survey_filters.QuestionFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		//Get All Questions
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Question, enums.Get)),
			middleware.ValidateDTO(survey_filters.QuestionFilters{}, middleware.QUERY, i),
		).Get("/all", h.GetAll)

		//Update Question
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Question, enums.Put)),
			middleware.ValidateDTO(survey_dtos.UpdateQuestionDTO{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.Update))

		//Delete Question
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Question, enums.Delete)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.Delete)
	})
}

func (h *QuestionHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, err := middleware.GetValidatedDTO[survey_dtos.CreateQuestionDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}
	result, createErr := h.questionService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *QuestionHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := r.URL.Query().Get("query")
	result, getErr := h.questionService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *QuestionHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
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

func (h *QuestionHandler) Update(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	dto, err := middleware.GetValidatedDTO[survey_dtos.UpdateQuestionDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.questionService.Update(ctx, objectId, dto)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *QuestionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	if deleteErr := h.questionService.Delete(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *QuestionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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
