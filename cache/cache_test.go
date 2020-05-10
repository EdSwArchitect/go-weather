package cache

import (
	"testing"

	"github.com/EdSwArchitect/go-weather/weather"
)

func TestCache(t *testing.T) {
	Initialize("http://localhost:9200")

	stations, err := weather.GetObservationStations()

	if err != nil {
		t.Errorf("Failed Getting stations. %+v\n", err)
	}

	InsertStations("stations", stations)

	//
	// PAJN
	//
	//
}

func TestContains(t *testing.T) {
	Initialize("http://localhost:9200")

	v := Contains("KCRG")

	if !v {
		t.Errorf("Should have found 'KCRG' in the index")
	}

	v = Contains("edwinfailed")

	if v {
		t.Errorf("Should NOT have found 'edwinfailed' in the index")
	}
}

func TestInsertFeatures(t *testing.T) {
	Initialize("http://localhost:9200")

	features, err := weather.GetFeatures()

	if err != nil {
		t.Errorf("Failed Getting features. %+v\n", err)
	}

	InsertFeatures("features", features)
}

func TestFeatureContains(t *testing.T) {
	Initialize("http://localhost:9200")

	v := ContainsFeature("KCRG")

	if !v {
		t.Errorf("Should have found 'KCRG' in the index")
	}

	v = ContainsFeature("edwinfailed")

	if v {
		t.Errorf("Should NOT have found 'edwinfailed' in the index")
	}
}
