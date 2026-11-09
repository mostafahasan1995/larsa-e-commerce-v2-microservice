package seed_location_data

import (
	"encoding/json"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
	"log"
	"os"
)

// Global variable to hold loaded area data
var Areas = []config_location_models.Area{}

// Intermediate structs to match the JSON structure from the 'areas.json' file

type AreaGeometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

type AreaProperties struct {
	NAME_1 string `json:"NAME_1"` // City/Governorate Name (e.g., AlḤasakah)
	NAME_2 string `json:"NAME_2"` // Area/District Name (e.g., Al-Hasakah)
}

type AreaFeature struct {
	Type       string         `json:"type"`
	Properties AreaProperties `json:"properties"`
	Geometry   AreaGeometry   `json:"geometry"`
}

type AreaFeatureCollection struct {
	Type     string        `json:"type"`
	Name     string        `json:"name"`
	Features []AreaFeature `json:"features"`
}

func InitAreas() []config_location_models.Area {
	filePath := "internal/modules/seed/data/json/areas.json"
	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading JSON file %s: %v", filePath, err)
		return []config_location_models.Area{}
	}

	// Unmarshal into a slice of FeatureCollections
	var featureCollections []AreaFeatureCollection
	err = json.Unmarshal(byteValue, &featureCollections)
	if err != nil {
		// Attempt to unmarshal as a single FeatureCollection as a fallback
		var singleFeatureCollection AreaFeatureCollection
		if errSingle := json.Unmarshal(byteValue, &singleFeatureCollection); errSingle == nil {
			featureCollections = []AreaFeatureCollection{singleFeatureCollection}
			log.Printf("Warning: areas.json was parsed as a single FeatureCollection, not an array.")
		} else {
			log.Printf("Error unmarshalling JSON data from %s (tried as array and single object): %v", filePath, err)
			return []config_location_models.Area{}
		}
	}

	var areas []config_location_models.Area

	// Iterate through each FeatureCollection in the slice
	for _, fc := range featureCollections {
		// Iterate through features within the current FeatureCollection
		for _, feature := range fc.Features {
			var multiPolygon [][][]float64
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

			seedArea := config_location_models.Area{
				BaseArea: config_location_models.BaseArea{
					Name:         feature.Properties.NAME_2,  // Use NAME_2 for Area/District Name
					CityName:     &feature.Properties.NAME_1, // Use NAME_1 for City/Governorate Name
					MultiPolygon: multiPolygon,
				},
			}
			areas = append(areas, seedArea)
		}
	}

	Areas = areas
	return areas
}
