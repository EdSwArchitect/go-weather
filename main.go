package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/EdSwArchitect/go-weather/cache"
	"github.com/EdSwArchitect/go-weather/weather"
	"github.com/gorilla/mux"
)

func init() {
	cache.Initialize("localhost:9200")
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
		cache.GetStationList("stations")
	}
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

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", heartBeat)
	router.HandleFunc("/stations", getStations)
	router.HandleFunc("/station/{stationId}", getStation)
	router.HandleFunc("/feature/{stationId}", getFeature)

	log.Fatal(http.ListenAndServe(":8080", router))
}
