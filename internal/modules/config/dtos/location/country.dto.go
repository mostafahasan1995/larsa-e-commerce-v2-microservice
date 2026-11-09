package config_location_dtos

type CreateCountryDto struct {
	Name          string `json:"name" validate:"required"`
	IsActive      *bool  `json:"isActive,omitempty"`
	Iso3          string `json:"iso3" validate:"required"`
	IncludeAreas  *bool  `json:"includeAreas,omitempty"`
	IncludeCities *bool  `json:"includeCities,omitempty"`
}
