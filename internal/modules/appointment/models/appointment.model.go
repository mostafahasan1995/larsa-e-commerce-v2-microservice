package appointment_models

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	appointment_enums "larsa-e-commerce-microservice-v2/internal/modules/appointment/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	profile_models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseAppointment struct {
	Date             time.Time                             `bson:"date" json:"date"`
	Time             string                                `bson:"time" json:"time"`
	TimeZone         string                                `bson:"timeZone" json:"timeZone"`
	Location         appointment_enums.AppointmentLocation `bson:"location" json:"location"`
	Customer         Customer                              `bson:"customer" json:"customer"`
	Status           appointment_enums.AppointmentStatus   `bson:"status" json:"status"`
	Notes            *Notes                                `bson:"notes,omitempty" json:"notes,omitempty"`
	UserCreateType   enums.UserType                        `bson:"userCreateType" json:"userCreateType"`
	models.BaseModel `bson:",inline"`
}

type Appointment struct {
	BaseAppointment `bson:",inline"`

	Translation *Translation        `bson:"translation,omitempty" json:"translation,omitempty"`
	VendorId    *primitive.ObjectID `bson:"vendorId" json:"vendorId,omitempty"`
	BuyerId     *primitive.ObjectID `bson:"buyerId" json:"buyerId,omitempty"`
}

type AppointmentResponse struct {
	BaseAppointment `bson:",inline"`

	Translation *TranslationResponse          `bson:"translation,omitempty" json:"translation,omitempty"`
	Vendor      *profile_models.VendorProfile `bson:"vendor,omitempty" json:"vendor,omitempty"`
	Buyer       *profile_models.BuyerProfile  `bson:"buyer,omitempty" json:"buyer,omitempty"`
}

type AppointmentsWithPagination = types.DataWithPagination[AppointmentResponse]

//

type Customer struct {
	Name         string `bson:"name" json:"name"`
	MobileNumber string `bson:"mobileNumber" json:"mobileNumber"`
	Email        string `bson:"email" json:"email"`
}

type Notes struct {
	Text       string              `bson:"text" json:"text"`
	Attachment *types.FileFieldDto `bson:"attachment" json:"attachment,omitempty"`
}
