package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

// ConvertToGroupModel converts a CreateGroupDto to a Group model
func ConvertToGroupModel(createGroupDto *dtos.CreateGroupDto, addionalOptions types.AdditionalOptions) *models.Group {
	group := &models.Group{
		BaseGroup: models.BaseGroup{
			Name:              createGroupDto.Name,
			Code:              *createGroupDto.Code,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(addionalOptions),
		},
	}
	return group
}
