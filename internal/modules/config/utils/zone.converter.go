package config_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_dto "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
)

func ConvertToZoneModel(dto *config_dto.CreateZoneDto, additionalOptions types.AdditionalOptions) config_models.Zone {
	zone := config_models.Zone{
		BaseZone: config_models.BaseZone{
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
			Name:              dto.Name,
			IsActive:          dto.IsActive,
			Code:              dto.Code,
			CountryIds:        dto.CountryIds,
			CityIds:           dto.CityIds,
			AreaIds:           dto.AreaIds,
			Polygons:          ConvertPolygonDtoToPolygon(dto.Polygons),
		},
	}

	return zone
}

func ConvertPolygonDtoToPolygon(dto []config_dto.PolygonDto) []config_models.Polygon {
	polygons := []config_models.Polygon{}
	for _, polygon := range dto {
		polygons = append(polygons, config_models.Polygon{
			Points: ConvertPointsDtoToPoints(polygon.Points),
		})
	}
	return polygons
}

func ConvertPointsDtoToPoints(dto []config_dto.PointDto) []config_models.Point {
	points := []config_models.Point{}
	for _, point := range dto {
		points = append(points, config_models.Point{
			Longitude: point.Longitude,
			Latitude:  point.Latitude,
		})
	}
	return points
}
