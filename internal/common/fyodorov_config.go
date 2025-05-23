package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/spf13/viper"
)

var DEFAULT_VERSION = "0.0.1"

type FyodorovConfig struct {
	Version               *string       `json:"version" yaml:"version,omitempty"`
	Providers             []Provider    `json:"providers,omitempty" yaml:"providers,omitempty"`
	Models                []ModelConfig `json:"models,omitempty" yaml:"models,omitempty"`
	Agents                []Agent       `json:"agents,omitempty" yaml:"agents,omitempty"`
	Tools                 []MCPTool     `json:"tools,omitempty" yaml:"tools,omitempty"`
	Instances             []Instance    `json:"instances,omitempty" yaml:"instances,omitempty"`
	TimeOfLastCacheUpdate *time.Time    `json:"time_of_last_cache_update,omitempty" yaml:"time_of_last_cache_update,omitempty"`
}

type Instance struct {
	ID      int64  `json:"id,omitempty" yaml:"id,omitempty"`
	Title   string `json:"title,omitempty" yaml:"title,omitempty"`
	AgentId int64  `json:"agent_id,omitempty" yaml:"agent_id,omitempty"`
}

func (i Instance) String() string {
	return fmt.Sprintf("%s-agent-%d", FormatString(i.Title), i.AgentId)
}

type BaseModel interface {
	String() string
	Validate() error
	GetID() int64
}
type Resource struct {
	ID int64 `json:"id,omitempty" yaml:"id,omitempty"`
}

func (r Resource) String() string {
	return fmt.Sprintf("resource-%d", r.ID)
}

func (r Resource) GetID() int64 {
	return r.ID
}

func (r Resource) Validate() error {
	return nil
}

func FormatString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ToValidUTF8(s, "")
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func CreateFyodorovConfig(v *viper.Viper) *FyodorovConfig {
	ttl := defaultTTL
	if v.IsSet("ttl") {
		ttl = v.GetDuration("ttl")
	}
	config := &FyodorovConfig{
		Version:   &DEFAULT_VERSION,
		Agents:    nil,
		Tools:     nil,
		Models:    nil,
		Providers: nil,
		Instances: nil,
	}
	t := time.Now().Add(ttl)
	config.TimeOfLastCacheUpdate = &t
	return config
}

func (c *FyodorovConfig) IsExpired(v *viper.Viper) bool {
	config, err := GetConfig(nil, v)
	if err != nil {
		fmt.Printf("Error getting config: %v\n", err)
		return true
	}
	return time.Since(*c.TimeOfLastCacheUpdate) > config.CacheTTL
}

func (config *FyodorovConfig) Validate() error {
	if config.Version != nil && *config.Version != "" {
		// check if version is in valid semver format
		if _, err := semver.NewVersion(*config.Version); err != nil {
			return err
		}
	}
	if config.Providers != nil {
		for _, provider := range config.Providers {
			if err := provider.Validate(); err != nil {
				return err
			}
		}
	}
	if config.Models != nil {
		for _, model := range config.Models {
			if err := model.Validate(); err != nil {
				return err
			}
		}
	}
	if config.Agents != nil {
		for _, agent := range config.Agents {
			if err := agent.Validate(); err != nil {
				return err
			}
		}
	}
	if config.Tools != nil {
		for _, tool := range config.Tools {
			if err := tool.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}
