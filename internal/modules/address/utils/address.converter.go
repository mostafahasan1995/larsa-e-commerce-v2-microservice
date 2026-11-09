package address_utils

import (
	address_dtos "larsa-e-commerce-microservice-v2/internal/modules/address/dtos"
	address_models "larsa-e-commerce-microservice-v2/internal/modules/address/models"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConvertAddressDTOToModel converts a CreateAddressDTO to a Address model
func ConvertAddressDTOToModel(dto address_dtos.AddressDTO, userId primitive.ObjectID) *address_models.AddressModel {

	now := time.Now()

	model := &address_models.AddressModel{
		Name:         &dto.Name,
		Type:         &dto.Type,
		Country:      &dto.Country,
		City:         &dto.City,
		Street:       &dto.Street,
		AddressLine1: &dto.AddressLine1,
		Notes:        dto.Notes,
		POBox:        dto.POBox,
		PostCode:     dto.PostCode,
		Phone:        &dto.Phone,
		IsDefault:    dto.IsDefault,
		UserId:       &userId,
		MetaDataBaseModel: models.MetaDataBaseModel{
			Id:        primitive.NewObjectID(),
			CreatedAt: now,
			CreatedBy: &userId,
			UpdatedAt: now,
			UpdatedBy: &userId,
		},
	}

	// Only set Location if it's not nil
	if dto.Location != nil {
		model.Location = &address_models.GeoLocation{
			Type:        dto.Location.Type,
			Coordinates: dto.Location.Coordinates,
		}
	}

	return model
}
