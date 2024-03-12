package common

import (
	"os"
	"path/filepath"
)

// Define a struct to hold the configuration
type Config struct {
	GagarinURL string `json:"gagarin_url"`
	// TsiolkovskyURL string `json:"tsiolkovsky_url"`
	// DostoyevskyURL string `json:"dostoyevsky_url"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
	Agents   []AgentClient `json:"agents"`
}

type AgentClient struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Instances []InstanceClient `json:"instances"`
}

type InstanceClient struct {
	ID    string `json:"id"`
	Title string `json:"name"`
}

func GetConfigPath() string {
	platform := os.Getenv("GOOS")
	switch platform {
	case "windows":
		return filepath.Join(GetPlatformBasePath(), "config.json")
	default:
		return filepath.Join(GetPlatformBasePath(), "config.json")
	}
}

func GetPlatformBasePath() string {
	platform := os.Getenv("GOOS")
	switch platform {
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "fyodorov")
	default:
		return filepath.Join(os.Getenv("HOME"), ".fyodorov")
	}
}
