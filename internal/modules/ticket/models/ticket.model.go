package ticket_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	payment_models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
	profile_models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	ticket_enums "larsa-e-commerce-microservice-v2/internal/modules/ticket/enums"
)

type BaseTicket struct {
	base_models.MetaDataBaseModel `bson:",inline"`

	Title       string                      `bson:"title" json:"title"`
	TypeId      primitive.ObjectID          `bson:"typeId" json:"typeId"`
	Status      ticket_enums.TicketStatus   `bson:"status" json:"status"`
	Priority    ticket_enums.TicketPriority `bson:"priority" json:"priority"`
	Description string                      `bson:"description" json:"description"`

	Attachment *types.FileFieldDto `bson:"attachment,omitempty" json:"attachment,omitempty"`

	QuotationId *primitive.ObjectID `bson:"quotationId,omitempty" json:"quotationId,omitempty"`
	InvoiceId   *primitive.ObjectID `bson:"invoiceId,omitempty" json:"invoiceId,omitempty"`
	OrderId     *primitive.ObjectID `bson:"orderId,omitempty" json:"orderId,omitempty"`

	RequesterId primitive.ObjectID  `bson:"requesterId" json:"requesterId"`
	BuyerId     *primitive.ObjectID `bson:"buyerId,omitempty" json:"buyerId,omitempty"`
	VendorId    *primitive.ObjectID `bson:"vendorId,omitempty" json:"vendorId,omitempty"`

	Comments []Comment `bson:"comments,omitempty" json:"comments,omitempty"`
}

type Ticket struct {
	BaseTicket `bson:",inline"`
}

type TicketResponse struct {
	BaseTicket `bson:",inline"`
	TicketType *config_models.TicketTypeResponse   `bson:"ticketType,omitempty" json:"ticketType,omitempty"`
	Requester  *profile_models.ProfileResponse     `bson:"requester,omitempty" json:"requester,omitempty"`
	Buyer      *profile_models.ProfileResponse     `bson:"buyer,omitempty" json:"buyer,omitempty"`
	Vendor     *profile_models.ProfileResponse     `bson:"vendor,omitempty" json:"vendor,omitempty"`
	Quotation  *quotation_models.QuotationResponse `bson:"quotation,omitempty" json:"quotation,omitempty"`
	Invoice    *payment_models.InvoiceResponse     `bson:"invoice,omitempty" json:"invoice,omitempty"`
	Order      *order_models.Order                 `bson:"order,omitempty" json:"order,omitempty"`
}

type Comment struct {
	Text       string              `bson:"text" json:"text"`
	Attachment *types.FileFieldDto `bson:"attachment,omitempty" json:"attachment,omitempty"`

	UserId    primitive.ObjectID          `bson:"userId" json:"userId"`
	User      *profile_models.BaseProfile `bson:"user,omitempty" json:"user,omitempty"`
	CreatedAt primitive.DateTime          `bson:"createdAt" json:"createdAt"`
}

type TicketsWithPagination = types.DataWithPagination[TicketResponse]
