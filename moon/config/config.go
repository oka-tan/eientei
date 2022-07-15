package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

//Config parametrizes Kaguya's configuration.
type Config struct {
	Boards         []string
	PostgresConfig PostgresConfig
	LnxConfig      LnxConfig
	LastModified   *time.Time
	BatchSize      uint
}

type PostgresConfig struct {
	ConnectionString string
}

type LnxConfig struct {
	Host string
	Port uint64
}

//LoadConfig reads config.json and unmarshals it into a Config struct.
//Errors might be returned due to IO or invalid JSON.
func LoadConfig() (Config, error) {
	blob, err := ioutil.ReadFile("./config.json")

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
