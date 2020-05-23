package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/EdSwArchitect/go-weather/weather"
	"github.com/cenkalti/backoff/v4"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

// Shards structure
type Shards struct {
	Total      int64 `json:"total"`
	Successful int64 `json:"successful"`
	Skipped    int64 `json:"skipped"`
	Failed     int64 `json:"failed"`
}

// CountResult structure
type CountResult struct {
	Count     int64  `json:"count"`
	ShardInfo Shards `json:"_shards"`
}

var es *elasticsearch.Client
var esHost string

// Initialize connection to ElasticSearch
func Initialize(host string) {

	esHost = host

	retryBackoff := backoff.NewExponentialBackOff()

	var err error

	es, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},

		RetryOnStatus: []int{502, 503, 504, 429},
		// Configure the backoff function
		//
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},

		// Retry up to 5 attempts
		//
		MaxRetries: 5,
	})

	if err != nil {
		log.Printf("Failed getting connection to elastic: %+v", err)
		return
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println(res)

}

// IndexCount get the index document count
func IndexCount(index string) (int64, error) {

	resp, err := http.Get(fmt.Sprintf("http://%s/%s/_count", esHost, index))

	if err != nil {
		return 0, err
	}

	var countResult CountResult

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&countResult)

	if err != nil {
		return 0, err
	}

	return countResult.Count, nil
}

// Contains - the station id is contained in the cache
func Contains(stationID string) bool {

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"_id": stationID,
			},
		},
	}

	err := json.NewEncoder(&buf).Encode(query)

	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := es.Count(
		es.Count.WithContext(context.Background()),
		es.Count.WithIndex("stations"),
		es.Count.WithBody(&buf),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	var e CountResult

	err = json.NewDecoder(res.Body).Decode(&e)

	if err != nil {
		log.Fatalf("Error decoding: %s", err)
	}

	log.Printf("The body: %+v", e)

	return e.Count == 1
}

// ContainsFeature - the station id is contained in the cache
func ContainsFeature(ID string) bool {

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"_id": ID,
			},
		},
	}

	err := json.NewEncoder(&buf).Encode(query)

	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := es.Count(
		es.Count.WithContext(context.Background()),
		es.Count.WithIndex("features"),
		es.Count.WithBody(&buf),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	var e CountResult

	err = json.NewDecoder(res.Body).Decode(&e)

	if err != nil {
		log.Fatalf("Error decoding: %s", err)
	}

	return e.Count == 1
}

// InsertStations inserts the stations into the Elastic index
func InsertStations(index string, stations weather.Stations) {

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        es,
		Index:         index,
		NumWorkers:    4,
		FlushBytes:    5e+6,
		FlushInterval: 30 * time.Second,
	})

	if err != nil {
		log.Fatalf("Unable to create bulk indexer: %s", err)
		return
	}

	var countSuccessful uint64

	for _, station := range stations.ObservationStations {

		var b strings.Builder
		b.WriteString(`{"station" : "`)
		b.WriteString(station)
		b.WriteString(`"}`)

		idx := strings.LastIndex(station, "/")

		stationRune := []rune(station)
		theID := string(stationRune[idx+1:])

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: theID,
				Body:       strings.NewReader(b.String()),

				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			log.Printf("Error bulk indexing: %s\n", err)
			return
		}
	} // for _, station := range stations.ObservationStations {

	if err = bi.Close(context.Background()); err != nil {

		log.Fatalf("Some fatal error: %s", err)
	}

	biStats := bi.Stats()

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	//
	log.Println(strings.Repeat("▔", 65))

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents",
			humanize.Comma(int64(biStats.NumFlushed)),
		)
	}

}

// InsertStations inserts the stations into the Elastic index
func InsertStationList(index string, stations []string) {

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        es,
		Index:         index,
		NumWorkers:    4,
		FlushBytes:    5e+6,
		FlushInterval: 30 * time.Second,
	})

	if err != nil {
		log.Fatalf("Unable to create bulk indexer: %s", err)
		return
	}

	var countSuccessful uint64

	for _, station := range stations {

		var b strings.Builder
		b.WriteString(`{"station" : "`)
		b.WriteString(station)
		b.WriteString(`"}`)

		idx := strings.LastIndex(station, "/")

		stationRune := []rune(station)
		theID := string(stationRune[idx+1:])

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: theID,
				Body:       strings.NewReader(b.String()),

				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			log.Printf("Error bulk indexing: %s\n", err)
			return
		}
	} // for _, station := range stations.ObservationStations {

	if err = bi.Close(context.Background()); err != nil {

		log.Fatalf("Some fatal error: %s", err)
	}

	biStats := bi.Stats()

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	//
	log.Println(strings.Repeat("▔", 65))

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents",
			humanize.Comma(int64(biStats.NumFlushed)),
		)
	}

}

// GetStationList inserts the stations into the Elastic index
func GetStationList(index string) ([]string, error) {
	// func GetStationList(index string) ([]map[string]string, error) {

	var buf bytes.Buffer

	// query := `{
	// 	"query": {
	// 		"match_all": {}
	// 	}
	// }`

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	err := json.NewEncoder(&buf).Encode(query)

	if err != nil {
		return nil, err
	}

	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			// Print the response status and error information.
			return nil, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])

		}
	}

	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// fmt.Printf("hits is: %v\n", r["hits"])

	// int(r["hits"])

	stations := make([]string, 1)
	// stations := make([]map[string]string, 1)

	//   // Print the response status, number of results, and request duration.
	//   log.Printf(
	// 	"[%s] %d hits; took: %dms",
	// 	res.Status(),
	// 	int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
	// 	int(r["took"].(float64)),
	//   )

	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		// stations[i] = fmt.Sprintf("%s", hit.(map[string]interface{})["_source"])

		// var m map[string]interface{}

		var m interface{}

		m = hit.(map[string]interface{})["_source"]

		switch dd := m.(type) {

		case map[string]interface{}:
			// fmt.Printf("map[string]interface{}: %s\n", dd)

			for _, v := range dd {
				// fmt.Printf("\t\tKey: %s. Value: %s  -- ", k, v)

				s, ok := v.(string)

				if ok {
					// fmt.Println("string")

					stations = append(stations, s)

				} else {
					_, ok = v.(interface{})
					if ok {
						fmt.Println(" and interface{}")
					}
				}

			}

		case map[string]string:
			// fmt.Printf("map[string]string: %s\n", dd)
		default:
			// fmt.Println("Something else")
		}

		// fmt.Printf("*** %v ***: %s\n-----\n", reflect.TypeOf(m), m)

		// stations = append(stations, fmt.Sprintf("%s", hit.(map[string]interface{})["_source"]))

	}

	// fmt.Printf("The stations returned: %+v", stations)

	return stations, nil
}

// InsertFeatures into the Elastic index
func InsertFeatures(index string, features []weather.Feature) {

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        es,
		Index:         index,
		NumWorkers:    4,
		FlushBytes:    5e+6,
		FlushInterval: 30 * time.Second,
	})

	if err != nil {
		log.Fatalf("Unable to create bulk indexer: %s", err)
	}

	var countSuccessful uint64

	for _, feature := range features {

		f, err := json.Marshal(feature)

		if err != nil {
			log.Fatalf("Unable to marshall object: %s", err)
		}

		var b strings.Builder
		b.WriteString(`{"feature" : `)
		b.WriteString(string(f))
		b.WriteString(`}`)

		idx := strings.LastIndex(feature.ID, "/")

		stationRune := []rune(feature.ID)

		theID := string(stationRune[idx+1:])

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: theID,
				Body:       strings.NewReader(b.String()),
				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			log.Printf("Error bulk indexing: %s\n", err)
			return
		}
	} // for _, station := range stations.ObservationStations {

	if err = bi.Close(context.Background()); err != nil {

		log.Fatalf("Some fatal error: %s", err)
	}

	biStats := bi.Stats()

	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	//
	log.Println(strings.Repeat("▔", 65))

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents",
			humanize.Comma(int64(biStats.NumFlushed)),
		)
	}

}
