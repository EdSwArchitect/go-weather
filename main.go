package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/EdSwArchitect/go-weather/weather"
	"github.com/gorilla/mux"
)

func heartBeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func getStations(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "getStations()")
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

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "%+v\n", feature)
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", heartBeat)
	router.HandleFunc("/stations", getStations)
	router.HandleFunc("/station/{stationId}", getStation)
	router.HandleFunc("/feature/{stationId}", getFeature)

	log.Fatal(http.ListenAndServe(":8080", router))
}
