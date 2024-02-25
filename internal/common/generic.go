package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Generic function to load configuration from a file.
// The format is determined by the file extension: .json or .yaml
func LoadConfig[T any](filename string) (*T, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error loading config from file %s: %v\n", filename, err)
		return nil, err
	}

	var config T
	switch filepath.Ext(filename) {
	case ".json":
		if err := json.Unmarshal(bytes, &config); err != nil {
			fmt.Printf("Error unmarshaling json config from file %s: %v\n", filename, err)
			return nil, err
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(bytes, &config); err != nil {
			fmt.Printf("Error unmarshaling yaml config from file %s: %v\n", filename, err)
			return nil, err
		}
	default:
		fmt.Printf("Error loading config from unsupported file format %s: %v\n", filename, err)
		return nil, fmt.Errorf("unsupported file format")
	}

	return &config, nil
}

// Define a generic type T which will be replaced by any type that is passed into SaveConfig
func SaveConfig[T any](config *T, filename string) error {
	var bytes []byte
	var err error
	switch filepath.Ext(filename) {
	case ".json":
		bytes, err = json.Marshal(config)
	case ".yaml", ".yml":
		bytes, err = yaml.Marshal(config)
	default:
		fmt.Printf("Unsupported file format %s: %v\n", filename, err)
		return fmt.Errorf("unsupported file format")
	}
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, 0644)
}
