package seed_location_data

import (
	"encoding/json"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
	"log"
	"os"
)

var Countries = []config_location_models.Country{}

// Intermediate structs to match the JSON structure
type Geometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

type Properties struct {
	GID_0   string `json:"GID_0"`
	COUNTRY string `json:"COUNTRY"`
}

type Feature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	Features []Feature `json:"features"`
}

func InitCountries() []config_location_models.Country {
	filePath := "internal/modules/seed/data/json/countries.json"
	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading JSON file %s: %v", filePath, err)
		return []config_location_models.Country{}
	}

	var featureCollections []FeatureCollection
	err = json.Unmarshal(byteValue, &featureCollections)
	if err != nil {
		log.Printf("Error unmarshalling JSON data from %s: %v", filePath, err)
		return []config_location_models.Country{}
	}

	var countries []config_location_models.Country

	for _, fc := range featureCollections {
		for _, feature := range fc.Features {
			var multiPolygon [][][]float64
			// JSON coordinates: [ MultiPolygon -> [ Polygon -> [ LinearRing -> [ Point -> [ lon, lat ] ] ] ] ]
			// Target structure: []Polygon -> []Point -> []float64
			for _, polygonCoords := range feature.Geometry.Coordinates {
				if len(polygonCoords) > 0 && len(polygonCoords[0]) > 0 {
					var polygonPoints [][]float64
					for _, pointCoords := range polygonCoords[0] {
						if len(pointCoords) == 2 {
							polygonPoints = append(polygonPoints, pointCoords)
						}
					}
					multiPolygon = append(multiPolygon, polygonPoints)
				}
			}

			seedCountry := config_location_models.Country{
				BaseCountry: config_location_models.BaseCountry{
					Name:         feature.Properties.COUNTRY,
					Iso3:         feature.Properties.GID_0,
					MultiPolygon: multiPolygon,
				},
			}
			countries = append(countries, seedCountry)
		}
	}

	Countries = countries
	return countries
}
