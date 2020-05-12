package weather

import (
	"encoding/json"
	"errors"
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
	ID           string    `json:"@id"`
	Type         string    `json:"@type"`
	TheElevation Elevation `json:"elevation"`
	StationID    string    `json:"stationIdentifier"`
	Name         string    `json:"name"`
	TimeZone     string    `json:"timeZone"`
}

// Weather type
type Weather struct {
	ID                  string     `json:"id"`
	Type                string     `json:"type"`
	Geo                 Geometry   `json:"geometry"`
	Props               Properties `json:"properties"`
	ObservationStations []string   `json:"observationStations"`
}

type Feature struct {
	ID    string     `json:"id"`
	Type  string     `json:"type"`
	Geo   Geometry   `json:"geometry"`
	Props Properties `json:"properties"`
}

type Stations struct {
	ObservationStations []string
}

type FeatureCollection struct {
	Type                string    `json:"type"`
	Features            []Feature `json:"features"`
	ObservationStations Stations  `json:"observationStations"`
}

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

func GetObservationStations() (Stations, error) {
	client := resty.New()

	resp, err := client.R().Get("https://api.weather.gov/stations")

	var stations Stations

	if err != nil {
		return stations, err
	}

	fmt.Printf("Status code of call: %d\n", resp.StatusCode())

	if resp.StatusCode() != 200 {
		return stations, errors.New(fmt.Sprintf("Status code returned: %d", resp.StatusCode()))
	}

	// fmt.Printf("%s\n", resp)

	err = json.Unmarshal([]byte(resp.String()), &stations)

	if err != nil {
		var s Stations
		fmt.Printf("Error parsing json\n")
		return s, err
	}

	// fmt.Printf("Size of observations: %d\n", len(w.ObservationStations))
	// fmt.Printf("The feature: %+v\n", w)

	return stations, nil

}

func GetStations() (string, error) {
	// https://api.weather.gov/stations

	client := resty.New()

	resp, err := client.R().Get("https://api.weather.gov/stations")

	if err != nil {
		return "", err
	}

	fmt.Printf("Status code of call: %d\n", resp.StatusCode())

	if resp.StatusCode() != 200 {
		return "", errors.New(fmt.Sprintf("Status code returned: %d", resp.StatusCode()))
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
		return nil, errors.New(fmt.Sprintf("Status code returned: %d", resp.StatusCode()))
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
