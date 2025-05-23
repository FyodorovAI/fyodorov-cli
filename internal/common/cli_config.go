package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	defaultTTL = 20 * time.Second
	JWT_TTL    = 20 * time.Minute
)

// Define a struct to hold the configuration
type Config struct {
	Version        string `json:"version"`
	GagarinURL     string `json:"gagarin_url"`
	TsiolkovskyURL string `json:"tsiolkovsky_url"`
	// DostoyevskyURL string `json:"dostoyevsky_url"`
	Email               string        `json:"email"`
	Password            string        `json:"password"`
	CacheTTL            time.Duration `json:"ttl"`
	JWT                 string        `json:"jwt"`
	TimeOfLastJWTUpdate time.Time     `json:"time_of_last_jwt_update"`
}

type AgentClient struct {
	ID        int64            `json:"id"`
	Name      string           `json:"name"`
	Instances []InstanceClient `json:"instances"`
}

type InstanceClient struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (c *Config) JWTExpired() bool {
	if c.JWT != "" && time.Since(c.TimeOfLastJWTUpdate) > JWT_TTL {
		return false
	}
	return true
}

func (c *Config) SetJWT(jwt string) error {
	c.JWT = jwt
	c.TimeOfLastJWTUpdate = time.Now()
	return SaveConfig(c, GetConfigPath())
}

func (c *Config) Validate() error {
	if c.Email == "" {
		return fmt.Errorf("email is required")
	}
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	if c.GagarinURL == "" {
		return fmt.Errorf("gagarin_url is required")
	}
	return nil
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

func GetConfig(config *Config, v *viper.Viper) (*Config, error) {
	if config == nil {
		config = &Config{}
	}
	if config.GagarinURL == "" {
		config.GagarinURL = v.GetString("gagarin-url")
	}
	if config.TsiolkovskyURL == "" {
		config.TsiolkovskyURL = v.GetString("tsiolkovsky-url")
	}
	if config.Email == "" {
		config.Email = v.GetString("email")
	}
	if config.Password == "" {
		config.Password = v.GetString("password")
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = defaultTTL
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	err := config.Validate()
	if err != nil {
		fmt.Printf("Error validating config: %v\n", err)
		return nil, err
	}
	return config, nil
}

func InitViper() *viper.Viper {
	v := viper.New()
	// Set default values
	v.SetDefault("gagarin-url", "https://gagarin.danielransom.com")
	v.SetDefault("tsiolkovsky-url", "https://tsiolkovsky.danielransom.com")
	// v.SetDefault("dostoyevsky-url", "https://dostoyevsky.danielransom.com")
	v.SetDefault("email", "")
	v.SetDefault("password", "")
	v.SetDefault("ttl", defaultTTL)
	v.SetDefault("jwt", "")
	v.SetDefault("jwt_ttl", JWT_TTL)
	v.SetDefault("time_of_last_jwt_update", time.Now().Add(-1*JWT_TTL))

	// Set the config file name and path
	v.SetConfigName("config")              // Name of the config file (without extension)
	v.SetConfigType("json")                // Config file format
	v.AddConfigPath(GetPlatformBasePath()) // Look for the config file in the user's home directory
	v.AddConfigPath(".")                   // Look for the config file in the current directory
	v.WatchConfig()                        // Watch for changes to the config file
	v.AutomaticEnv()                       // Read in environment variables that match
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Printf("No config file found in local directory or at %s\n", GetPlatformBasePath())
			return v
		} else {
			// Config file was found but another error was produced
			fmt.Printf("Error reading config file: %v\n", err)
			return v
		}
	}
	return v
}
