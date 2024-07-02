package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"database"`
	// Add other configuration fields as needed
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalf("Error unmarshalling config data: %v", err)
		return nil, err
	}

	return config, nil
}
