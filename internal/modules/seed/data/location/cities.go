package seed_location_data

import (
	"encoding/json"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
	"log"
	"os"
)

var Cities = []config_location_models.City{}

type CityGeometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

type CityProperties struct {
	COUNTRY string `json:"COUNTRY"` // Country Name
	NAME_1  string `json:"NAME_1"`  // City/Region Name (e.g., AlHasakah)
}

type CityFeature struct {
	Type       string         `json:"type"`
	Properties CityProperties `json:"properties"`
	Geometry   CityGeometry   `json:"geometry"`
}

type CityFeatureCollection struct {
	Type     string        `json:"type"`
	Name     string        `json:"name"`
	Features []CityFeature `json:"features"`
}

func InitCities() []config_location_models.City {
	filePath := "internal/modules/seed/data/json/cities.json"
	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading JSON file %s: %v", filePath, err)
		return []config_location_models.City{}
	}

	// Unmarshal into a slice of FeatureCollections
	var featureCollections []CityFeatureCollection
	err = json.Unmarshal(byteValue, &featureCollections)
	if err != nil {
		// Attempt to unmarshal as a single FeatureCollection as a fallback
		var singleFeatureCollection CityFeatureCollection
		if errSingle := json.Unmarshal(byteValue, &singleFeatureCollection); errSingle == nil {
			featureCollections = []CityFeatureCollection{singleFeatureCollection}
			log.Printf("Warning: cities.json was parsed as a single FeatureCollection, not an array.")
		} else {
			log.Printf("Error unmarshalling JSON data from %s (tried as array and single object): %v", filePath, err)
			return []config_location_models.City{}
		}
	}

	var cities []config_location_models.City

	// Iterate through each FeatureCollection in the slice
	for _, fc := range featureCollections {
		// Iterate through features within the current FeatureCollection
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

			seedCity := config_location_models.City{
				BaseCity: config_location_models.BaseCity{
					Name:         feature.Properties.NAME_1, // Use NAME_1 for city/region name
					CountryName:  &feature.Properties.COUNTRY,
					MultiPolygon: multiPolygon,
				},
			}
			cities = append(cities, seedCity)
		}
	}

	Cities = cities
	return cities
}
