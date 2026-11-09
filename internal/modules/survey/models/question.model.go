package survey_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	survey_enums "larsa-e-commerce-microservice-v2/internal/modules/survey/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionModel struct {
	Type                     survey_enums.QuestionType   `bson:"type" json:"type"`
	Order                    float64                     `bson:"order" json:"order"`
	Question                 types.Localizable[string]   `bson:"question" json:"question"`
	Options                  types.Localizable[[]string] `bson:"options" json:"options"`
	VendorId                 *primitive.ObjectID         `bson:"vendorId,omitempty" json:"vendorId,omitempty"`
	BranchIds                []primitive.ObjectID        `bson:"branchIds" json:"branchIds"`
	AllBranches              bool                        `bson:"allBranches" json:"allBranches"`
	models.MetaDataBaseModel `bson:",inline"`
}

type QuestionResponse struct {
	QuestionModel `bson:",inline"`
	Branches      []*product_models.BranchResponse `bson:"branches,omitempty" json:"branches,omitempty"`
}

type QuestionsWithPagination = types.DataWithPagination[QuestionResponse]
