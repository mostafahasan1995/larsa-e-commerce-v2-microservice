package appointment_handler

import (
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	appointment_dto "larsa-e-commerce-microservice-v2/internal/modules/appointment/dto"
	appointment_filters "larsa-e-commerce-microservice-v2/internal/modules/appointment/filters"
	appointment_service "larsa-e-commerce-microservice-v2/internal/modules/appointment/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentHandler struct {
	appointmentService appointment_service.AppointmentService
	validationInstance *validator.Validate
}

func NewAppointmentHandler(i *do.Injector, r *chi.Mux) {
	h := AppointmentHandler{
		appointmentService: do.MustInvoke[appointment_service.AppointmentService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.APPOINTMENT.GetAdminRoute(), func(r chi.Router) {
		// Create Appointment
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Appointment, enums.Post)),
			middleware.ValidateDTO(appointment_dto.CreateAppointmentDTO{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		// Get Appointment By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Appointment, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		// Get Appointments With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Appointment, enums.Get)),
			middleware.ValidateDTO(appointment_filters.AppointmentFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		// Reschedule Appointment
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Appointment, enums.Put)),
			middleware.ValidateDTO(appointment_dto.RescheduleAppointmentDTO{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}/reschedule", helpers.Make(h.Reschedule))

		// Request Translation
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Appointment, enums.Put)),
			middleware.ValidateDTO(appointment_dto.TranslationRequestDTO{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}/translation", helpers.Make(h.RequestTranslation))

		// Update Partial
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Appointment, enums.Put)),
			middleware.ValidateDTO(appointment_dto.UpdateAppointmentPartialDTO{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Patch("/{id}/partial", helpers.Make(h.UpdatePartial))

	})
}

func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, _ := middleware.GetValidatedDTO[appointment_dto.CreateAppointmentDTO](r)

	result, createErr := h.appointmentService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *AppointmentHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, getErr := h.appointmentService.GetById(ctx, objectId)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *AppointmentHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.appointmentService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *AppointmentHandler) Reschedule(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, _ := middleware.GetValidatedDTO[appointment_dto.RescheduleAppointmentDTO](r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, rescheduleErr := h.appointmentService.Reschedule(ctx, objectId, dto)
	if rescheduleErr != nil {
		helpers.HandleApiError(w, rescheduleErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *AppointmentHandler) RequestTranslation(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, _ := middleware.GetValidatedDTO[appointment_dto.TranslationRequestDTO](r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, translationErr := h.appointmentService.RequestTranslation(ctx, objectId, dto)
	if translationErr != nil {
		helpers.HandleApiError(w, translationErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *AppointmentHandler) UpdatePartial(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, _ := middleware.GetValidatedDTO[appointment_dto.UpdateAppointmentPartialDTO](r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, updateErr := h.appointmentService.UpdatePartial(ctx, objectId, dto)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}
