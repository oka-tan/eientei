package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	PostgresConfig PostgresConfig
	Hosts          []string
	Production     bool
}

type PostgresConfig struct {
	ConnectionString string
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
