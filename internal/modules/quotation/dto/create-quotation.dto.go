package quotation_dto

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateQuotationDTO struct {
	Title    string                        `json:"title" validate:"required"`
	Type     quotation_enums.QuotationType `json:"type" `
	Deadline *time.Time                    `json:"deadline"`

	Products []QuotationProductDTO `json:"products" validate:"required,min=1"`

	Delivery       DeliveryDTO     `json:"delivery" validate:"required"`
	Payment        PaymentDTO      `json:"payment" validate:"required"`
	DiscountAmount *QuotationPrice `json:"discountAmount"`

	WarrantyAndValidity WarrantyAndValidityDTO `json:"warrantyAndValidity" validate:"required"`

	PackagingRequirements PackagingRequirementsDTO `json:"packagingRequirements" validate:"required"`

	AdditionalNotes *AdditionalNotesDTO  `json:"additionalNotes"`
	Attachments     []types.FileFieldDto `json:"attachments"`
	Signature       *SignatureDTO        `json:"signature"`

	QuotationBuyerInfo *QuotationBuyerInfoDTO `json:"quotationBuyerInfo" bson:"quotationBuyerInfo"`

	VendorProfileId *primitive.ObjectID `json:"vendorProfileId,omitempty"`

	ProductsObjects  []product_models.ProductResponse    `json:"-"`
	RequestQuotation *quotation_models.QuotationResponse `json:"-"`
	//TODO: update postman docs.
}

func (dto *CreateQuotationDTO) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	productRepo := base_repos.NewBaseRepo[product_models.Product, product_models.ProductResponse](args.DB, enums.PRODUCT)

	productsIds := make([]primitive.ObjectID, len(dto.Products))
	for i, product := range dto.Products {
		productsIds[i] = product.ProductId
	}

	products, err := productRepo.GetByIds(args.Context, productsIds, []bson.M{}, base_utils.AdditionalOptions{})
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid products",
		}
	}

	// Validate that all requested products were found
	if len(products) != len(productsIds) {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Some requested products were not found",
		}
	}

	dto.ProductsObjects = products

	return nil, nil
}

//TODO: ask ms. Ruba about:
//1- buyer website notes
//2- vendor Add Custom quote to customer
