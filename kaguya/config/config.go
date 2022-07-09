package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Config parametrizes Kaguya's configuration.
type Config struct {
	APIConfig        APIConfig
	ImagesConfig     ImagesConfig
	ThumbnailsConfig ImagesConfig
	Boards           []BoardConfig
	PostgresConfig   PostgresConfig
	InitialNap       string
	SkipArchive      bool
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
	S3Host         string
	Region         string
	Bucket         string
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
	Thumbnails  bool
	Images      bool
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

func (c *Config) StartImagesService() bool {
	for _, b := range c.Boards {
		if b.Images {
			return true
		}
	}

	return false
}

func (c *Config) StartThumbnailsService() bool {
	for _, b := range c.Boards {
		if b.Thumbnails {
			return true
		}
	}

	return false
}
