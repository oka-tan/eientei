package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Config parametrizes Kaguya's configuration.
type Config struct {
	APIConfig      APIConfig
	ImagesConfig   ImagesConfig
	Boards         []BoardConfig
	PostgresConfig PostgresConfig
}

//APIConfig parametrizes Kaguya's configuration for the consumption of 4chan's API.
type APIConfig struct {
	RequestTimeout string
	Host           string
}

//ImagesConfig parametrizes Kaguya's configuration for downloading media from 4chan and posting it to S3.
type ImagesConfig struct {
	NapTime        string
	RequestTimeout string
	Host           string
	AwsRegion      string
	BucketName     string
}

//PostgresConfig parametrizes Kaguya's configuration for PostgreSQL
type PostgresConfig struct {
	ConnectionString string
}

//BoardConfig parametrizes Kaguya's configuration for each board being scraped
type BoardConfig struct {
	Name        string
	NapTime     string
	LongNapTime string
}

//LoadConfig reads config.json and unmarshals it into a Config struct.
//Errors might be returned due to IO or invalid JSON.
func LoadConfig() (Config, error) {
	configFile := os.Args[1]
	blob, err := ioutil.ReadFile(configFile)

	if err != nil {
		return Config{}, fmt.Errorf("Error loading file config.json in project root: %s", err)
	}

	var conf Config

	err = json.Unmarshal(blob, &conf)

	if err != nil {
		return Config{}, fmt.Errorf(
			"Error unmarshalling configuration file contents to JSON:\n File contents: %s\n Error message: %s",
			blob,
			err,
		)
	}

	return conf, nil
}
