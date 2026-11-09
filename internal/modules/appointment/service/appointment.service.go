package appointment_service

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	appointment_dto "larsa-e-commerce-microservice-v2/internal/modules/appointment/dto"
	appointment_enums "larsa-e-commerce-microservice-v2/internal/modules/appointment/enums"
	appointment_filters "larsa-e-commerce-microservice-v2/internal/modules/appointment/filters"
	appointment_models "larsa-e-commerce-microservice-v2/internal/modules/appointment/models"
	appointment_utils "larsa-e-commerce-microservice-v2/internal/modules/appointment/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentRepo interface {
	Create(ctx context.Context, appointment appointment_models.Appointment) (*appointment_models.AppointmentResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID) (*appointment_models.AppointmentResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter appointment_filters.AppointmentFilters) (*appointment_models.AppointmentsWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]appointment_models.AppointmentResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, appointment appointment_models.Appointment, filter bson.M) (*appointment_models.AppointmentResponse, *helpers.ApiError)
	UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}, userId primitive.ObjectID) (*appointment_models.AppointmentResponse, *helpers.ApiError)
}

type AppointmentService interface {
	Create(ctx context.Context, dto *appointment_dto.CreateAppointmentDTO) (*appointment_models.AppointmentResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID) (*appointment_models.AppointmentResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (appointment_models.AppointmentsWithPagination, *helpers.ApiError)
	Reschedule(ctx context.Context, id primitive.ObjectID, dto *appointment_dto.RescheduleAppointmentDTO) (*appointment_models.AppointmentResponse, *helpers.ApiError)
	RequestTranslation(ctx context.Context, id primitive.ObjectID, dto *appointment_dto.TranslationRequestDTO) (*appointment_models.AppointmentResponse, *helpers.ApiError)
	UpdatePartial(ctx context.Context, id primitive.ObjectID, dto *appointment_dto.UpdateAppointmentPartialDTO) (*appointment_models.AppointmentResponse, *helpers.ApiError)
}

type appointment struct {
	repo AppointmentRepo
}

func NewAppointmentService(i *do.Injector) (AppointmentService, error) {
	return &appointment{
		repo: do.MustInvoke[AppointmentRepo](i),
	}, nil
}

func (s *appointment) Create(ctx context.Context, dto *appointment_dto.CreateAppointmentDTO) (*appointment_models.AppointmentResponse, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get user information",
		}
	}
	appointment, apiErr := appointment_utils.ConvertCreateAppointmentDTOToModel(dto, ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	appointment.Status = appointment_enums.PENDING
	if cfg.User.Type == enums.Vendor {
		appointment.UserCreateType = enums.Vendor
		appointment.VendorId = &cfg.User.Id
	} else {
		appointment.UserCreateType = enums.Buyer
		appointment.BuyerId = &cfg.User.Id
	}

	result, apiErr := s.repo.Create(ctx, appointment)
	if apiErr != nil {
		return nil, apiErr
	}
	return result, nil
}

func (s *appointment) Reschedule(ctx context.Context, id primitive.ObjectID, dto *appointment_dto.RescheduleAppointmentDTO) (*appointment_models.AppointmentResponse, *helpers.ApiError) {
	_, apiErr := s.GetById(ctx, id)
	if apiErr != nil {
		return nil, apiErr
	}
	cfg, _ := util.GetReqAppCfg(ctx)

	updateFields := make(map[string]interface{})
	updateFields["date"] = dto.Date
	updateFields["time"] = dto.Time

	result, apiErr := s.repo.UpdateFieldsById(ctx, id, updateFields, cfg.User.Id)
	if apiErr != nil {
		return nil, apiErr
	}
	return result, nil
}

func (s *appointment) RequestTranslation(ctx context.Context, id primitive.ObjectID, dto *appointment_dto.TranslationRequestDTO) (*appointment_models.AppointmentResponse, *helpers.ApiError) {
	_, apiErr := s.GetById(ctx, id)
	if apiErr != nil {
		return nil, apiErr
	}
	cfg, _ := util.GetReqAppCfg(ctx)

	updateFields := make(map[string]interface{})
	updateFields["translation"] = appointment_models.Translation{
		BaseTranslation: appointment_models.BaseTranslation{
			TranslateType: dto.TranslateType,
		},
		TranslateFromId: dto.TranslateFromId,
		TranslateToId:   dto.TranslateToId,
		InvoiceId:       dto.InvoiceId,
	}

	return s.repo.UpdateFieldsById(ctx, id, updateFields, cfg.User.Id)
}

func (s *appointment) GetById(ctx context.Context, id primitive.ObjectID) (*appointment_models.AppointmentResponse, *helpers.ApiError) {
	return s.repo.GetById(ctx, id)
}

func (s *appointment) GetWithPagination(ctx context.Context, skip, limit int64, query string) (appointment_models.AppointmentsWithPagination, *helpers.ApiError) {
	filter, err := appointment_filters.NewAppointmentFilters(query)
	if err != nil {
		return appointment_models.AppointmentsWithPagination{}, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	result, apiErr := s.repo.GetWithPagination(ctx, skip, limit, *filter)
	if apiErr != nil {
		return appointment_models.AppointmentsWithPagination{}, apiErr
	}

	return *result, nil
}

func (s *appointment) UpdatePartial(ctx context.Context, id primitive.ObjectID, dto *appointment_dto.UpdateAppointmentPartialDTO) (*appointment_models.AppointmentResponse, *helpers.ApiError) {
	appointment, apiErr := s.GetById(ctx, id)
	if apiErr != nil {
		return nil, apiErr
	}

	cfg, _ := util.GetReqAppCfg(ctx)

	updateFields := make(map[string]interface{})

	if dto.Status != nil {
		apiErr := validateUpdateStatus(*appointment, *dto.Status, cfg.User.Id)
		if apiErr != nil {
			return nil, apiErr
		}
		updateFields["status"] = *dto.Status
	}
	if dto.Customer != nil {
		updateFields["customer"] = appointment_models.Customer{
			Name:         dto.Customer.Name,
			MobileNumber: dto.Customer.MobileNumber,
			Email:        dto.Customer.Email,
		}
	}
	if dto.Notes != nil {
		updateFields["notes"] = appointment_models.Notes{
			Text:       dto.Notes.Text,
			Attachment: dto.Notes.Attachment,
		}
	}

	//TODO: Send Notification to the vendor or buyer
	return s.repo.UpdateFieldsById(ctx, id, updateFields, cfg.User.Id)
}

func validateUpdateStatus(appointment appointment_models.AppointmentResponse, status appointment_enums.AppointmentStatus, userId primitive.ObjectID) *helpers.ApiError {
	previousStatus := appointment.Status

	if previousStatus == appointment_enums.PENDING {
		if status == appointment_enums.APPROVED || status == appointment_enums.REJECTED {
			if appointment.UserCreateType == enums.Buyer {
				if appointment.Vendor.Id == userId {
					return &helpers.ApiError{
						StatusCode: 400,
						Message:    "Only the guest can approve or reject a pending appointment",
					}
				}
			} else {
				if appointment.Buyer.Id == userId {
					return &helpers.ApiError{
						StatusCode: 400,
						Message:    "Only the guest can approve or reject a pending appointment",
					}
				}
			}
		}
	}

	return nil
}
