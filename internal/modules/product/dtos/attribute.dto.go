package product_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OptionDto struct {
	Label string `json:"label" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type CreateAttributeDto struct {
	Name           types.Localizable[string]         `json:"name" validate:"required,localizable"`
	Code           string                            `json:"code" validate:"required"`
	Options        types.LocalizableArray[OptionDto] `json:"options" validate:"omitempty,min=0,dive,required"`
	GroupId        primitive.ObjectID                `json:"groupId" validate:"required"`
	RefField       *product_enums.AttributeRefField  `json:"refField" validate:"omitempty,oneof= branch addon"`
	Type           product_enums.AttributeType       `json:"type" validate:"required,oneof=text number email password date time datetime_local radio checkbox file color range tel url search select textarea wysiwyg_editor switch custom"`
	Multi          bool                              `json:"multi"`
	RelatedTo      *[]string                         `json:"relatedTo" validate:"omitempty,min=1,dive,required"`
	IsFiltrable    bool                              `json:"isFiltrable"`
	IsComparable   bool                              `json:"isComparable"`
	Visibility     bool                              `json:"visibility"`
	Required       bool                              `json:"required"`
	Sortable       bool                              `json:"sortable"`
	SelectionStyle *string                           `json:"selectionStyle" validate:"required"`
	BelongsTo      []product_enums.BelongTo          `json:"belongsTo" validate:"required,min=1,dive,oneof=product quotation production_line"`
	IsLocalizable  bool                              `json:"isLocalizable"`
	Configurable   bool                              `json:"configurable"`
}

func (dto *CreateAttributeDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	groupRepo := base_repos.NewBaseRepo[product_models.Group, product_models.GroupResponse](args.DB, enums.GROUP)
	attributeRepo := base_repos.NewBaseRepo[product_models.Attribute, product_models.AttributeResponse](args.DB, enums.ATTRIBUTE)

	if dto.Type != product_enums.CUSTOM_GROUP && len(dto.Options) == 0 && dto.RefField == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Options is required for non-group type",
		}
	}

	if dto.Type == product_enums.CUSTOM_GROUP && dto.RelatedTo == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "RelatedTo is required for group type",
		}
	} else if dto.Type != product_enums.CUSTOM_GROUP && dto.RelatedTo != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "RelatedTo is not allowed for non-group type",
		}
	}

	if dto.Type == product_enums.CUSTOM_GROUP && dto.RelatedTo != nil {
		if dto.IsLocalizable {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "IsLocalizable is not allowed for custom group type",
			}
		}
		attributes, err := attributeRepo.GetAll(args.Context, []bson.M{
			{
				"$match": bson.M{
					"code": bson.M{
						"$in": *dto.RelatedTo,
					},
					"tag": bson.M{
						"$eq": product_enums.CUSTOM,
					},
					"type": bson.M{
						"$ne": product_enums.CUSTOM_GROUP,
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}

		if len(attributes) != len(*dto.RelatedTo) {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "Some of the RelatedTo attributes are not found",
			}
		}
	}

	_, err := groupRepo.GetById(args.Context, dto.GroupId, []bson.M{})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
