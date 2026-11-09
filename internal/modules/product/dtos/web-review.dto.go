package product_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
)

type CreateReviewDto struct {
	Rate    uint8  `json:"rate" validate:"required,min=0,max=5"`
	Comment string `json:"comment" validate:"required,min=1,max=1000"`
}

func (dto *CreateReviewDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {

	return nil, nil
}
