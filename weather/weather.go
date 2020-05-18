package weather

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

// Geometry type
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// Elevation type
type Elevation struct {
	Value    float64 `json:"value"`
	UnitCode string  `json:"unitCode"`
}

// Properties type
type Properties struct {
	ID              string    `json:"@id"`
	Type            string    `json:"@type"`
	TheElevation    Elevation `json:"elevation"`
	StationID       string    `json:"stationIdentifier"`
	Name            string    `json:"name"`
	TimeZone        string    `json:"timeZone"`
	Forecast        string    `json:"forecast"`
	County          string    `json:"county"`
	FireWeatherZone string    `json:"fireWeatherZone"`
}

// Weather type
type Weather struct {
	ID                  string     `json:"id"`
	Type                string     `json:"type"`
	Geo                 Geometry   `json:"geometry"`
	Props               Properties `json:"properties"`
	ObservationStations []string   `json:"observationStations"`
}

// Feature The weather feature
type Feature struct {
	ID    string     `json:"id"`
	Type  string     `json:"type"`
	Geo   Geometry   `json:"geometry"`
	Props Properties `json:"properties"`
}

// Stations The stations
type Stations struct {
	ObservationStations []string
}

// FeatureCollection the feature collection
type FeatureCollection struct {
	Type                string    `json:"type"`
	Features            []Feature `json:"features"`
	ObservationStations Stations  `json:"observationStations"`
}

// FB a type....
type FB struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// ParseWeather function
func ParseWeather(jayson string) (Feature, error) {

	fmt.Printf("Parsing feature:\n%s\n----\n", jayson)
	var feature Feature

	err := json.Unmarshal([]byte(jayson), &feature)

	if err != nil {
		return feature, err
	}

	return feature, nil

}

// GetObservationStations ....
func GetObservationStations() (Stations, error) {
	client := resty.New()

	resp, err := client.R().Get("https://api.weather.gov/stations")

	var stations Stations

	if err != nil {

		log.Printf("Getting the stations list failed: %+v", err)

		return stations, err
	}

	fmt.Printf("Status code of call: %d\n", resp.StatusCode())

	if resp.StatusCode() != 200 {

		return stations, fmt.Errorf("Status code returned: %d", resp.StatusCode())
	}

	// fmt.Printf("%s\n", resp)

	err = json.Unmarshal([]byte(resp.String()), &stations)

	if err != nil {
		var s Stations
		fmt.Printf("Error parsing json\n")
		return s, err
	}

	fmt.Printf("Size of observations: %d\n", len(stations.ObservationStations))

	return stations, nil
}

// GetStations get the stations
func GetStations() (string, error) {
	// https://api.weather.gov/stations

	client := resty.New()

	resp, err := client.R().Get("https://api.weather.gov/stations")

	if err != nil {
		return "", err
	}

	fmt.Printf("Status code of call: %d\n", resp.StatusCode())

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("Status code returned: %d", resp.StatusCode())
	}

	fmt.Printf("%s\n", resp)

	var w Feature

	w, err = ParseWeather(resp.String())

	if err != nil {
		fmt.Printf("Error parsing json\n")
	}

	fmt.Printf("The id: %+s\n", w.ID)

	// fmt.Printf("Size of observations: %d\n", len(w.ObservationStations))
	// fmt.Printf("The feature: %+v\n", w)

	return "Success", nil
}

func GetFeatures() ([]Feature, error) {
	// https://api.weather.gov/stations

	client := resty.New()

	resp, err := client.R().Get("https://api.weather.gov/stations")

	if err != nil {
		return nil, err
	}

	// fmt.Printf("Status code of call: %d\n", resp.StatusCode())

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code returned: %d", resp.StatusCode())
	}

	// fmt.Printf("%s\n", resp)

	var features FB

	err = json.Unmarshal([]byte(resp.String()), &features)

	if err != nil {
		log.Printf("Failed unmarshalling into features %s", err)
		return nil, err
	}

	return features.Features, nil
}

// GetFeature for the station ID
func GetFeature(stationID string) (Feature, error) {
	// https://api.weather.gov/stations

	client := resty.New()

	resp, err := client.R().Get(fmt.Sprintf("https://api.weather.gov/stations/%s", stationID))

	if err != nil {
		return Feature{}, err
	}

	// fmt.Printf("Status code of call: %d\n", resp.StatusCode())

	if resp.StatusCode() != 200 {
		return Feature{}, fmt.Errorf("Status code returned: %d", resp.StatusCode())
	}

	// fmt.Printf("%s\n", resp)

	var feature Feature

	err = json.Unmarshal([]byte(resp.String()), &feature)

	if err != nil {
		log.Printf("Failed unmarshalling into features %s", err)
		return Feature{}, err
	}

	return feature, nil
}
