package geocoding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	nominatimReverseURL = "https://nominatim.openstreetmap.org/reverse"
	nominatimSearchURL  = "https://nominatim.openstreetmap.org/search"
	userAgent           = "LarsaECommerce/1.0"
	defaultTimeout      = 5 * time.Second
)

// nominatimResponse represents the JSON response from Nominatim reverse API.
type nominatimResponse struct {
	Address map[string]interface{} `json:"address"`
}

// nominatimSearchResult represents a single result from Nominatim search API.
type nominatimSearchResult struct {
	Lat     string                 `json:"lat"`
	Lon     string                 `json:"lon"`
	Address map[string]interface{} `json:"address"`
}

// ReverseGeocode converts (lat, lng) to (country, city, area).
// Inputs: latitude, longitude
// Outputs: country, city, area (empty string if not found)
func ReverseGeocode(lat, lng float64) (country, city, area string, err error) {
	if os.Getenv("TEST_MODE") == "true" || os.Getenv("GO_ENV") == "test" {
		return "", "", "", nil
	}

	u, err := url.Parse(nominatimReverseURL)
	if err != nil {
		return "", "", "", err
	}
	q := u.Query()
	q.Set("format", "json")
	q.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Set("lon", strconv.FormatFloat(lng, 'f', -1, 64))
	q.Set("addressdetails", "1")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("nominatim returned status %d", resp.StatusCode)
	}

	var data nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", "", err
	}

	if data.Address == nil {
		return "", "", "", nil
	}

	country = getStr(data.Address, "country")
	// For rural/sparse areas (e.g. Libya, Chad), city/town/village may be empty; use state/county/region as fallbacks
	city = getFirst(data.Address, "city", "town", "village", "municipality", "state_district", "state", "county", "region")
	// Area = more local level; avoid duplicating state in both city and area
	area = getFirst(data.Address, "suburb", "neighbourhood", "quarter", "district", "locality")

	return country, city, area, nil
}

// Geocode converts an address string to coordinates and location components.
// Returns lat, lng, country, city, area (empty strings if not found).
func Geocode(address string) (lat, lng float64, country, city, area string, err error) {
	if os.Getenv("TEST_MODE") == "true" || os.Getenv("GO_ENV") == "test" {
		return 0, 0, "", "", "", nil
	}

	if strings.TrimSpace(address) == "" {
		return 0, 0, "", "", "", fmt.Errorf("address cannot be empty")
	}

	u, parseErr := url.Parse(nominatimSearchURL)
	if parseErr != nil {
		return 0, 0, "", "", "", parseErr
	}
	q := u.Query()
	q.Set("q", address)
	q.Set("format", "json")
	q.Set("limit", "1")
	q.Set("addressdetails", "1")
	u.RawQuery = q.Encode()

	req, reqErr := http.NewRequest(http.MethodGet, u.String(), nil)
	if reqErr != nil {
		return 0, 0, "", "", "", reqErr
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: defaultTimeout}
	resp, respErr := client.Do(req)
	if respErr != nil {
		return 0, 0, "", "", "", respErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, "", "", "", fmt.Errorf("nominatim search returned status %d", resp.StatusCode)
	}

	var results []nominatimSearchResult
	if decodeErr := json.NewDecoder(resp.Body).Decode(&results); decodeErr != nil {
		return 0, 0, "", "", "", decodeErr
	}

	if len(results) == 0 {
		return 0, 0, "", "", "", fmt.Errorf("no results found for address: %s", address)
	}

	result := results[0]

	// Parse latitude and longitude
	lat, latErr := strconv.ParseFloat(result.Lat, 64)
	if latErr != nil {
		return 0, 0, "", "", "", fmt.Errorf("invalid latitude in response: %s", result.Lat)
	}

	lng, lngErr := strconv.ParseFloat(result.Lon, 64)
	if lngErr != nil {
		return 0, 0, "", "", "", fmt.Errorf("invalid longitude in response: %s", result.Lon)
	}

	// Extract location components using same logic as ReverseGeocode
	if result.Address != nil {
		country = getStr(result.Address, "country")
		city = getFirst(result.Address, "city", "town", "village", "municipality", "state_district", "state", "county", "region")
		area = getFirst(result.Address, "suburb", "neighbourhood", "quarter", "district", "locality")
	}

	return lat, lng, country, city, area, nil
}

func getStr(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok || strings.TrimSpace(s) == "" {
		return ""
	}
	return s
}

func getFirst(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if s := getStr(m, key); s != "" {
			return s
		}
	}
	return ""
}
