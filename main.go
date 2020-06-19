package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/EdSwArchitect/go-weather/cache"
	"github.com/EdSwArchitect/go-weather/config"
	"github.com/EdSwArchitect/go-weather/weather"
	"github.com/gorilla/mux"
)

var espUri string
var configFile string
var featuresURI string
var stationsURI string
var httpPort int

func init() {

	flag.StringVar(&espUri, "espUri", "localhost:9200", "The ESP host and port number")
	flag.IntVar(&httpPort, "serverPort", 8080, "The HTTP server port")
	flag.StringVar(&configFile, "configFile", "", "The configuration file")

	flag.Parse()

	if configFile != "" {
		log.Printf("Working with log file: %s", configFile)

		config, err := config.ReadConfig(&configFile)

		if err != nil {
			log.Fatalf("Configuration failed: %s", err)
		}

		espUri = config.EspURI
		featuresURI = config.FeaturesURI
		stationsURI = config.StationsURI
		httpPort = config.ServerPort

	}

	log.Printf("espURI: %s", espUri)
	log.Printf("configFile: %s", configFile)
	log.Printf("featuresURI: %s", featuresURI)
	log.Printf("stationsURI: %s", stationsURI)
	log.Printf("httpPort: %d", httpPort)

	// cache.Initialize("localhost:9200")
	cache.Initialize(espUri)
}

func heartBeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println("OK")
}

func getStations(w http.ResponseWriter, r *http.Request) {

	count, err := cache.IndexCount("stations")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	if count == 0 {

		theStations, err := weather.GetObservationStations()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("content-type", "text/plain; charset=utf-8")
			fmt.Fprintf(w, "%s\n", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("content-type", "application/json; charset=utf-8")

		fmt.Fprintf(w, "%s", theStations.ObservationStations)
	} else {
		//		cache.GetStations()
		stations, err := cache.GetStationList("stations")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Add("content-type", "text/plain; charset=utf-8")
			fmt.Fprintf(w, "%s\n", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("content-type", "application/json; charset=utf-8")

		for _, station := range stations {
			fmt.Fprintf(w, "%s\n", station)
		}

		// fmt.Fprintf(w, "%s", stations)

	}
}

func loadStations(w http.ResponseWriter, r *http.Request) {

	theStations, err := weather.GetObservationStations()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	cache.InsertStationList(stationsURI, theStations.ObservationStations)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "text/plain; charset=utf-8")

	fmt.Fprint(w, "OK")
}

func getStation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	stationID := vars["stationId"]

	if stationID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "No stationId given")
		return
	}

	feature, err := weather.GetFeature(stationID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "No stationId %s found", stationID)
		return
	}

	b, err := json.Marshal(feature)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Unable to marshal station information")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json; charset=utf-8")

	fmt.Fprintf(w, "%s", string(b))
}

func loadFeatures(w http.ResponseWriter, r *http.Request) {

	features, err := weather.GetFeatures()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Unable to marshal feature information. %s", err)
		return
	}

	cache.InsertFeatures(featuresURI, features)
}

func getFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Fprintf(w, "The vars: %+v", vars)

	stationID := vars["stationId"]

	if stationID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "No stationId given")
		return
	}

	feature, err := weather.GetFeature(stationID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "No stationId %s found", stationID)
		return
	}

	b, err := json.Marshal(feature)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Unable to marshal station information")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json; charset=utf-8")

	fmt.Fprintf(w, "%s", string(b))
}

func getFeatures(w http.ResponseWriter, r *http.Request) {

	// count, err := cache.IndexCount("features")

	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Header().Add("content-type", "text/plain; charset=utf-8")
	// 	fmt.Fprintf(w, "%s\n", err)
	// 	return
	// }

	// if count == 0 {

	theFeatures, err := weather.GetFeatures()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json; charset=utf-8")

	fmt.Fprintf(w, "%v", theFeatures)
	// } else {
	//		cache.GetStations()
	// 	stations, err := cache.

	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Header().Add("content-type", "text/plain; charset=utf-8")
	// 		fmt.Fprintf(w, "%s\n", err)
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusOK)
	// 	w.Header().Add("content-type", "application/json; charset=utf-8")

	// 	for _, station := range stations {
	// 		fmt.Fprintf(w, "%s\n", station)
	// 	}

	// 	// fmt.Fprintf(w, "%s", stations)

	// }
}

func writeStatic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Fprintf(w, "The vars: %+v", vars)

	staticID := vars["staticID"]

	if staticID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "No staticID given")
		return
	}

	file, err := os.Create(fmt.Sprintf("/perm-data/%s", staticID))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Unable to write file: /perm-data/%s", staticID)
		return
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s: %s", time.Now(), staticID))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Add("content-type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Unable to write content to file /perm-data/%s", staticID)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "OK")
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", heartBeat)
	router.HandleFunc("/stations", getStations)
	router.HandleFunc("/features", getFeatures)
	router.HandleFunc("/loadStations", loadStations)
	router.HandleFunc("/station/{stationId}", getStation)
	router.HandleFunc("/loadFeatures", loadFeatures)
	router.HandleFunc("/feature/{stationId}", getFeature)
	router.HandleFunc("/writeStatic/{saticID}", writeStatic)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), router))
}
