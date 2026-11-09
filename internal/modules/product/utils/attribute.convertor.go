package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

func convertOptions(dtoOptions types.LocalizableArray[product_dtos.OptionDto]) types.LocalizableArray[product_models.Option] {
	options := make(types.LocalizableArray[product_models.Option], len(dtoOptions))
	for i, opt := range dtoOptions {
		response := make(types.Localizable[product_models.Option])
		for key, value := range opt {
			response[enums.LanguageKey(key)] = product_models.Option{
				Label: value.Label,
				Value: value.Value,
			}
		}
		options[i] = response
	}
	return options
}

// ConvertToAttributeModel converts a CreateAttributeDto to an Attribute model
func ConvertToAttributeModel(createAttributeDto *product_dtos.CreateAttributeDto, additionalOptions types.AdditionalOptions) *product_models.Attribute {

	attribute := &product_models.Attribute{
		BaseAttribute: product_models.BaseAttribute{
			VendorId: &additionalOptions.UserId,

			Name:              createAttributeDto.Name,
			Options:           convertOptions(createAttributeDto.Options),
			Type:              createAttributeDto.Type,
			BelongsTo:         createAttributeDto.BelongsTo,
			GroupId:           &createAttributeDto.GroupId,
			Tag:               product_enums.CUSTOM,
			IsFiltrable:       createAttributeDto.IsFiltrable,
			Required:          createAttributeDto.Required,
			IsComparable:      createAttributeDto.IsComparable,
			Visibility:        createAttributeDto.Visibility,
			Sortable:          createAttributeDto.Sortable,
			SelectionStyle:    createAttributeDto.SelectionStyle,
			Multi:             createAttributeDto.Multi,
			RefField:          createAttributeDto.RefField,
			Code:              createAttributeDto.Code,
			RelatedTo:         createAttributeDto.RelatedTo,
			IsLocalizable:     createAttributeDto.IsLocalizable,
			Configurable:      createAttributeDto.Configurable,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
	return attribute
}

//
