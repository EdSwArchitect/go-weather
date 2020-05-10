package weather

import (
	"fmt"
	"testing"
)

func TestParsing(t *testing.T) {
	fmt.Println("Running TestParsing")

	ex := `{
		"id": "https://api.weather.gov/stations/KSFO",
		"type": "Feature",
		"geometry": {
		"type": "Point",
		"coordinates": [
		-122.36558,
		37.61961
		]
		},
		"properties": {
		"@id": "https://api.weather.gov/stations/KSFO",
		"@type": "wx:ObservationStation",
		"elevation": {
		"value": 3.048,
		"unitCode": "unit:m"
		},
		"stationIdentifier": "KSFO",
		"name": "San Francisco, San Francisco International Airport",
		"timeZone": "America/Los_Angeles"
		}
		}`

	j, err := ParseWeather(ex)

	if err != nil {
		t.Errorf("Failed parsing the JSON. %+v\n", err)
	}

	if j.ID != "https://api.weather.gov/stations/KSFO" {
		t.Errorf("ID not as expected. %s\n", j.ID)
	}

	if j.Type != "Feature" {
		t.Errorf("Type is not Feature. %s\n", j.Type)
	}

	if j.Geo.Type != "Point" {
		t.Errorf("j.Geo.Type != Point: %s\n", j.Geo.Type)
	}

	if j.Geo.Coordinates[0] != -122.36558 || j.Geo.Coordinates[1] != 37.61961 {
		t.Errorf("Coordinates are off: %+v\n", j.Geo.Coordinates)
	}

	fmt.Printf("Parsed object: %+v\n", j)
}

func TestUrlCall(t *testing.T) {
	ans, err := GetObservationStations()

	if err != nil {
		t.Errorf("Getting stations failed: %+v\n", err)
	}

	fmt.Printf("The results size is: %d\n", len(ans.ObservationStations))
}
