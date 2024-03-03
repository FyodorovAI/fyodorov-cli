package common

import "fmt"

type ModelInfo struct {
	Mode               string   `json:"mode" yaml:"mode"`
	InputCostPerToken  *float64 `json:"input_cost_per_token" yaml:"input_cost_per_token"`
	OutputCostPerToken *float64 `json:"output_cost_per_token" yaml:"output_cost_per_token"`
	MaxTokens          *int     `json:"max_tokens" yaml:"max_tokens"`
	BaseModel          string   `json:"base_model" yaml:"base_model"`
}

type Params struct{}

type ModelConfig struct {
	Name      string     `json:"name" yaml:"name"`
	Provider  string     `json:"provider" yaml:"provider"`
	Params    Params     `json:"params" yaml:"params"`
	ModelInfo *ModelInfo `json:"model_info" yaml:"model_info"`
}

var MODEL_MODES = []string{"embedding", "chat", "multimodal"}

func (c *ModelInfo) Validate() error {
	if !contains(MODEL_MODES, c.Mode) {
		return fmt.Errorf("invalid model mode: %s", c.Mode)
	}
	if c.InputCostPerToken != nil && *c.InputCostPerToken < 0.0 {
		return fmt.Errorf("invalid input cost per token: %f", *c.InputCostPerToken)
	}
	if c.OutputCostPerToken != nil && *c.OutputCostPerToken < 0.0 {
		return fmt.Errorf("invalid output cost per token: %f", *c.OutputCostPerToken)
	}
	if c.MaxTokens != nil && *c.MaxTokens < 2 {
		return fmt.Errorf("invalid max tokens: %d", c.MaxTokens)
	}
	return nil
}

func (c *ModelConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("model name is required")
	}
	if c.Provider == "" {
		return fmt.Errorf("model provider is required")
	}
	return c.ModelInfo.Validate()
}

func contains(array []string, element string) bool {
	for _, item := range array {
		if item == element {
			return true
		}
	}
	return false
}
