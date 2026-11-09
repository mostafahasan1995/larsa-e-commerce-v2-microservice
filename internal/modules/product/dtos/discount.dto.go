package product_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HandleNewProductsDiscountDto struct {
	ProductIds []primitive.ObjectID `json:"productIds" validate:"required,min=1,dive,required"`
	Discount   DiscountDto          `json:"discount" validate:"required"`
	BranchId   *primitive.ObjectID  `json:"branchId" validate:"omitempty"`
	IsDefault  bool                 `json:"isDefault" validate:"omitempty"`
}

func (dto *HandleNewProductsDiscountDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	productRepo := base_repos.NewBaseRepo[product_models.Product, product_models.ProductResponse](args.DB, enums.PRODUCT)
	branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.DB, enums.BRANCH)
	if dto.BranchId == nil && !dto.IsDefault {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "You must set IsDefault to true if you don't set a branchId",
		}
	}

	if dto.BranchId != nil && dto.IsDefault {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "You must set IsDefault to false if you set a branchId",
		}
	}

	if dto.BranchId != nil {
		_, err := branchRepo.GetById(args.Context, *dto.BranchId, []bson.M{
			{"$match": base_utils.AddSoftDeletedFilter(bson.M{})},
		})
		if err != nil {
			return nil, err
		}
	}

	_, err := productRepo.GetByIds(args.Context, dto.ProductIds, []bson.M{
		{"$match": base_utils.AddSoftDeletedFilter(bson.M{})},
	}, base_utils.AdditionalOptions{ReturnError: true})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
