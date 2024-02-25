package common

import (
	"encoding/json"
	"os"
)

// Define a struct to hold the configuration
type Config struct {
	GagarinURL     string `json:"gagarin_url"`
	TsiolkovskyURL string `json:"tsiolkovsky_url"`
	DostoyevskyURL string `json:"dostoyevsky_url"`
	Email          string `json:"email"`
	Password       string `json:"password"`
}

func (config *Config) LoadConfig(path string) {
	file, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(file, config)
	}
}

func (config *Config) SaveConfig(path string) {
	file, _ := json.MarshalIndent(config, "", " ")
	os.WriteFile(path, file, 0644)
}
