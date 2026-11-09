package appointment_dto

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	appointment_enums "larsa-e-commerce-microservice-v2/internal/modules/appointment/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerDTO struct {
	Name         string `json:"name" validate:"required"`
	MobileNumber string `json:"mobileNumber" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
}

type NotesDTO struct {
	Text       string              `json:"text"`
	Attachment *types.FileFieldDto `json:"attachment"`
}

type TranslationRequestDTO struct {
	TranslateFromId *primitive.ObjectID             `json:"translateFromId" validate:"required"`
	TranslateToId   *primitive.ObjectID             `json:"translateToId" validate:"required"`
	TranslateType   appointment_enums.TranslateType `json:"translateType" validate:"required"`
	InvoiceId       *primitive.ObjectID             `json:"invoiceId" validate:"required"`
}

type CreateAppointmentDTO struct {
	Date     time.Time                             `json:"date" validate:"required"`
	Time     string                                `json:"time" validate:"required"`
	TimeZone string                                `json:"timeZone" validate:"required"`
	Location appointment_enums.AppointmentLocation `json:"location" validate:"required"`

	Customer CustomerDTO `json:"customer" validate:"required"`

	Notes *NotesDTO `json:"notes"`

	VendorId *primitive.ObjectID `json:"vendorId"`
	BuyerId  *primitive.ObjectID `json:"buyerId"`
}

type RescheduleAppointmentDTO struct {
	Date time.Time `json:"date" validate:"required"`
	Time string    `json:"time" validate:"required"`
}

type UpdateAppointmentPartialDTO struct {
	Notes    *NotesDTO                            `json:"notes,omitempty"`
	Customer *CustomerDTO                         `json:"customer,omitempty"`
	Status   *appointment_enums.AppointmentStatus `json:"status,omitempty"`
}
