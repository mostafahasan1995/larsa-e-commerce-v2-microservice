package survey_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	survey_dtos "larsa-e-commerce-microservice-v2/internal/modules/survey/dtos"
	survey_services "larsa-e-commerce-microservice-v2/internal/modules/survey/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type WebFeedbackHandler struct {
	feedbackService    survey_services.WebFeedbackService
	validationInstance *validator.Validate
}

func NewWebFeedbackHandler(i *do.Injector, r *chi.Mux) {
	h := WebFeedbackHandler{
		feedbackService:    do.MustInvoke[survey_services.WebFeedbackService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.FEEDBACK.GetWebRoute(), func(r chi.Router) {
		// Submit feedback
		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(survey_dtos.CreateFeedbackDTO{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))
	})
}

func (h *WebFeedbackHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[survey_dtos.CreateFeedbackDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.feedbackService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}
