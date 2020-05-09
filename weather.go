package main

import (
	"encoding/json"
	"errors"
	"fmt"

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

type W struct {
	Id string `json:"id"`
}

// ParseWeather function
func ParseWeather(jayson string) (Weather, error) {
	var weather Weather

	err := json.Unmarshal([]byte(jayson), &weather)

	if err != nil {
		return weather, err
	}

	return weather, nil

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

	var w Weather

	w, err = ParseWeather(resp.String())

	if err != nil {
		fmt.Printf("Error parsing json\n")
	}

	fmt.Printf("Size of observations: %d\n", len(w.ObservationStations))

	return "Success", nil
}
