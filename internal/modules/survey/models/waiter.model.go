package survey_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WaiterModel struct {
	Name                     types.Localizable[string] `bson:"name" json:"name"`
	Image                    *types.FileFieldDto       `bson:"image" json:"image"`
	BranchId                 primitive.ObjectID        `bson:"branchId" json:"branchId"`
	VendorId                 *primitive.ObjectID       `bson:"vendorId,omitempty" json:"vendorId,omitempty"`
	models.MetaDataBaseModel `bson:",inline"`
}

type WaiterResponse struct {
	WaiterModel `bson:",inline"`
	Branch      *product_models.BranchResponse `bson:"branch,omitempty" json:"branch,omitempty"`
}

type WaitersWithPagination = types.DataWithPagination[WaiterResponse]
