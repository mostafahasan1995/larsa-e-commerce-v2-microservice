package appointment_utils

import (
	"context"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	appointment_dto "larsa-e-commerce-microservice-v2/internal/modules/appointment/dto"
	appointment_models "larsa-e-commerce-microservice-v2/internal/modules/appointment/models"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConvertCreateAppointmentDTOToModel converts CreateAppointmentDTO to Appointment model
func ConvertCreateAppointmentDTOToModel(dto *appointment_dto.CreateAppointmentDTO, ctx context.Context) (appointment_models.Appointment, *helpers.ApiError) {
	now := time.Now()

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return appointment_models.Appointment{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	appointment := appointment_models.Appointment{
		BaseAppointment: appointment_models.BaseAppointment{
			BaseModel: models.BaseModel{
				Id: primitive.NewObjectID(),
				MetaDataBaseModel: models.MetaDataBaseModel{
					CreatedAt: now,
					CreatedBy: &cfg.User.Id,
					UpdatedAt: now,
					UpdatedBy: &cfg.User.Id,
				},
			},
			Date:     dto.Date,
			Time:     dto.Time,
			TimeZone: dto.TimeZone,
			Location: dto.Location,
			Customer: appointment_models.Customer{
				Name:         dto.Customer.Name,
				MobileNumber: dto.Customer.MobileNumber,
				Email:        dto.Customer.Email,
			},
		},
	}

	if dto.Notes != nil {
		appointment.Notes = &appointment_models.Notes{
			Text:       dto.Notes.Text,
			Attachment: dto.Notes.Attachment,
		}
	}

	return appointment, nil
}
