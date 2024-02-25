package common

import (
	"fmt"
)

type AgentConfig struct {
	ProviderID          string   `json:"provider_id" yaml:"provider_id"`
	Model               string   `json:"model" yaml:"model"`
	NameForHuman        string   `json:"name_for_human" yaml:"name_for_human"`
	DescriptionForHuman string   `json:"description_for_human" yaml:"description_for_human"`
	Prompt              string   `json:"prompt" yaml:"prompt"`
	PromptSize          int      `json:"prompt_size" yaml:"prompt_size"`
	Tools               []string `json:"tools" yaml:"tools"`
	Rag                 []string `json:"rag" yaml:"rag"`
}

// The format is determined by the file extension: .json or .yaml
func LoadAgentConfigFromFile(filename string) (*AgentConfig, error) {
	config, err := LoadConfig[AgentConfig](filename)
	if err != nil {
		fmt.Printf("Error loading agent config from file %s: %v\n", filename, err)
		return nil, err
	}

	return config, nil
}

// A function to save an AgentConfig to either a yaml or json file
func (config *AgentConfig) SaveToFile(filename string) error {
	SaveConfig[AgentConfig](config, filename)

	return nil
}
