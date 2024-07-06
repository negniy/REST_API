package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port string `json:"port"`
}

func LoadConfig() (*Config, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
