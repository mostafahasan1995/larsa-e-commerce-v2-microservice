package config_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
)

// ConvertToTicketTypeModel converts a CreateTicketTypeDto to a TicketType model
func ConvertToTicketTypeModel(dto *config_dtos.CreateTicketTypeDto, additionalOptions types.AdditionalOptions) *config_models.TicketType {
	ticketType := &config_models.TicketType{
		BaseTicketType: config_models.BaseTicketType{
			Code:                 dto.Code,
			Name:                 dto.Name,
			ExpectedSolveTime:    dto.ExpectedSolveTime,
			UpdateNotification:   dto.UpdateNotification,
			InternalNotification: dto.InternalNotification,
			TeamNotifyIds:        dto.TeamNotifyIds,
			IsActive:             dto.IsActive,
			MetaDataBaseModel:    util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
	return ticketType
}
