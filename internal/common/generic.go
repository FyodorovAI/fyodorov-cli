package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigLoader is an interface for loading configuration from a file.
type BaseConfig interface {
	Validate() error
}

// Generic function to load configuration from a file.
// The format is determined by the file extension: .json or .yaml
func LoadConfig[T any](filename string) (*T, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Expand environment variables in the file contents
	//  Any ${FOO} will be replaced with os.Getenv("FOO")
	expanded := os.ExpandEnv(string(fileBytes))

	var config T
	switch filepath.Ext(filename) {
	case ".json":
		if err := json.Unmarshal([]byte(expanded), &config); err != nil {
			fmt.Printf("\033[33mError unmarshaling json config from file %s:\n  %v\n\033[0m", filename, err)
			return nil, err
		}
	case ".yaml", ".yml":
		dec := yaml.NewDecoder(bytes.NewReader([]byte(expanded)))
		dec.KnownFields(true) // ← reject any unknown fields
		if err := dec.Decode(&config); err != nil {
			fmt.Printf("\033[33mError unmarshaling yaml config from file %s:\n  %v\n\033[0m", filename, err)
			return nil, err
		}
	default:
		fmt.Printf("\033[33mError loading config from unsupported file format %s:\n  %v\n\033[0m", filename, err)
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
		bytes, err = json.MarshalIndent(config, "", "  ")
	case ".yaml", ".yml":
		bytes, err = yaml.Marshal(config)
	default:
		fmt.Printf("\033[33mUnsupported file format %s: %v\n\033[0m", filename, err)
		return fmt.Errorf("unsupported file format")
	}
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, 0644)
}
