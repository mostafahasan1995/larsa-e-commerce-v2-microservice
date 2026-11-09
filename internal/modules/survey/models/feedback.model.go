package survey_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Answer struct {
	QuestionId primitive.ObjectID `bson:"questionId" json:"questionId"`
	Answer     string             `bson:"answer" json:"answer"`
}

type FeedbackModel struct {
	Phone                    string              `bson:"phone" json:"phone"`
	ClientName               string              `bson:"clientName" json:"clientName"`
	Email                    string              `bson:"email" json:"email"`
	Birthday                 string              `bson:"birthday" json:"birthday"`
	WaiterId                 primitive.ObjectID  `bson:"waiterId" json:"waiterId"`
	BranchId                 *primitive.ObjectID `bson:"branchId" json:"branchId"`
	VendorId                 primitive.ObjectID  `bson:"vendorId" json:"vendorId"`
	Points                   float64             `bson:"points" json:"points"`
	Answers                  []Answer            `bson:"answers" json:"answers"`
	Notes                    string              `bson:"notes" json:"notes"`
	WaiterRate               float64             `bson:"waiterRate" json:"waiterRate"`
	models.MetaDataBaseModel `bson:",inline"`
}

type FeedbackResponse struct {
	FeedbackModel `bson:",inline"`
	Waiter        *WaiterResponse                `bson:"waiter,omitempty" json:"waiter,omitempty"`
	Branch        *product_models.BranchResponse `bson:"branch,omitempty" json:"branch,omitempty"`
	Questions     []QuestionResponse             `bson:"questions,omitempty" json:"questions,omitempty"`
	CouponCode    *string                        `bson:"couponCode,omitempty" json:"couponCode,omitempty"`
}

type FeedbacksWithPagination = types.DataWithPagination[FeedbackResponse]
