package survey_handlers

import (
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	survey_filters "larsa-e-commerce-microservice-v2/internal/modules/survey/filters"
	survey_services "larsa-e-commerce-microservice-v2/internal/modules/survey/services"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedbackHandler struct {
	feedbackService    survey_services.FeedbackService
	validationInstance *validator.Validate
}

func NewFeedbackHandler(i *do.Injector, r *chi.Mux) {
	h := FeedbackHandler{
		feedbackService:    do.MustInvoke[survey_services.FeedbackService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.FEEDBACK.GetAdminRoute(), func(r chi.Router) {

		//Get Feedback By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feedback, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		//Get Feedbacks With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feedback, enums.Get)),
			middleware.ValidateDTO(survey_filters.FeedbackFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		//Get All Feedbacks
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feedback, enums.Get)),
			middleware.ValidateDTO(survey_filters.FeedbackFilters{}, middleware.QUERY, i),
		).Get("/all", h.GetAll)

		//Export Feedbacks
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feedback, enums.Get)),
		).Get("/export", h.Export)

		//Delete Feedback
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feedback, enums.Delete)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.Delete)
	})
}

func (h *FeedbackHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	query := r.URL.Query().Get("query")
	result, getErr := h.feedbackService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *FeedbackHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.feedbackService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *FeedbackHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	if deleteErr := h.feedbackService.Delete(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FeedbackHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, getErr := h.feedbackService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *FeedbackHandler) Export(w http.ResponseWriter, r *http.Request) {
	ctx, _ := util.AddCtxAppCfg(r)

	excelData, exportErr := h.feedbackService.ExportFeedbacks(ctx)
	if exportErr != nil {
		helpers.HandleApiError(w, exportErr)
		return
	}
	fileName := fmt.Sprintf("feedbacks_export_%s.xlsx", time.Now().Format("2006-01-02_15-04-05"))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Write(excelData)
}
