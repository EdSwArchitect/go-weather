package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

/*
{
	"espUri" : "localhost:9200",
	"stationsIndex" : "stations",
	"serverPort" : 80,
	"featuresIndex"	 : "features"
}
*/

type Config struct {
	EspURI      string `json:"espUri"`
	StationsURI string `json:"stationsIndex"`
	ServerPort  int    `json:"serverPort"`
	FeaturesURI string `json:"featuresIndex"`
}

// ReadConfig read the configuraion file
func ReadConfig(path *string) (Config, error) {
	if *path == "" {
		return Config{}, fmt.Errorf("Unable to read configuration file")
	}

	jsonFile, err := os.Open(*path)

	if err != nil {
		return Config{}, err
	}

	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(bytes, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}
