package common

import (
	"github.com/Masterminds/semver"
)

var DEFAULT_VERSION = "0.0.1"

type FyodorovConfig struct {
	Version   string        `json:"version" yaml:"version,omitempty"`
	Providers []Provider    `json:"providers,omitempty" yaml:"providers,omitempty"`
	Models    []ModelConfig `json:"models,omitempty" yaml:"models,omitempty"`
	Agents    []Agent       `json:"agents,omitempty" yaml:"agents,omitempty"`
	Tools     []MCPTool     `json:"tools,omitempty" yaml:"tools,omitempty"`
}

func CreateFyodorovConfig() *FyodorovConfig {
	return &FyodorovConfig{
		Version: DEFAULT_VERSION,
		Agents:  nil,
		Tools:   nil,
	}
}

func (config *FyodorovConfig) Validate() error {
	if config.Version != "" {
		// check if version is in valid semver format
		if _, err := semver.NewVersion(config.Version); err != nil {
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
