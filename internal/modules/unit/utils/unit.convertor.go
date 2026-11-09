package unit_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/unit/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/unit/models"
)

// ConvertToUnitModel converts a CreateUnitDto to a Unit model
func ConvertToUnitModel(createUnitDto *dtos.CreateUnitDto, addionalOptions types.AdditionalOptions) *models.Unit {

	unit := &models.Unit{
		BaseUnit: models.BaseUnit{
			Name:              createUnitDto.Name,
			Type:              createUnitDto.Type,
			Category:          createUnitDto.Category,
			Rate:              createUnitDto.Rate,
			Active:            createUnitDto.Active,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(addionalOptions),
		},
	}
	return unit
}
