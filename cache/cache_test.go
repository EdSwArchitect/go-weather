package cache

import (
	"log"
	"testing"

	"github.com/EdSwArchitect/go-weather/weather"
)

func init() {
	log.Printf("init")
}

// TestCache test caching
func TestCache(t *testing.T) {

	log.Println("TestCache")

	Initialize("http://localhost:9200")

	stations, err := weather.GetObservationStations()

	if err != nil {
		t.Errorf("Failed Getting stations. %+v\n", err)
	}

	log.Printf("****** Stations length: %d", len(stations.ObservationStations))

	InsertStations("stations", stations)

	//
	// PAJN
	//
	//
}

func TestContains(t *testing.T) {

	log.Println("TestContains")

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

	log.Println("TestInsertFeatures")

	Initialize("http://localhost:9200")

	features, err := weather.GetFeatures()

	if err != nil {
		t.Errorf("Failed Getting features. %+v\n", err)
	}

	InsertFeatures("features", features)
}

func TestFeatureContains(t *testing.T) {

	log.Println("TestFeatureContains")

	Initialize("http://localhost:9200")

	v := ContainsFeature("KEFK")

	if !v {
		t.Errorf("Should have found 'KEFK' in the index")
	}

	v = ContainsFeature("edwinfailed")

	if v {
		t.Errorf("Should NOT have found 'edwinfailed' in the index")
	}
}

func TestStationsCache(t *testing.T) {
	stations, err := GetStationList("stations")

	if err != nil {
		t.Errorf("Failed getting list of stations from cache: %+v\n", err)
	}

	log.Printf("Stations size: %d\n%v\n", len(stations), stations)
}
