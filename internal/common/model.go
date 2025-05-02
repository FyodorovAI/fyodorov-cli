package common

import "fmt"

type ModelInfo struct {
	Mode               string   `json:"mode" yaml:"mode,omitempty"`
	InputCostPerToken  *float64 `json:"input_cost_per_token" yaml:"input_cost_per_token,omitempty"`
	OutputCostPerToken *float64 `json:"output_cost_per_token" yaml:"output_cost_per_token,omitempty"`
	MaxTokens          *int     `json:"max_tokens" yaml:"max_tokens,omitempty"`
	BaseModel          string   `json:"base_model" yaml:"base_model,omitempty"`
}

type Params struct{}

type ModelConfig struct {
	ID        int64      `json:"id,omitempty" yaml:"id,omitempty"`
	Name      string     `json:"name,omitempty" yaml:"name,omitempty"`
	Provider  string     `json:"provider,omitempty" yaml:"provider,omitempty"`
	Params    Params     `json:"params,omitempty" yaml:"params,omitempty"`
	ModelInfo *ModelInfo `json:"model_info" yaml:"model_info,omitempty"`
}

var MODEL_MODES = Enum{"embedding", "chat", "multimodal"}

func (c *ModelInfo) Validate() error {
	if c.BaseModel == "" {
		return fmt.Errorf("base model is required")
	}
	if c.Mode != "" && !MODEL_MODES.Contains(c.Mode) {
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
	if c.ModelInfo != nil {
		if err := c.ModelInfo.Validate(); err != nil {
			return fmt.Errorf("model info is invalid: %v", err)
		}
	} else {
		return fmt.Errorf("model info is required")
	}
	return nil
}

func (c *ModelConfig) GetModelHandle() string {
	return c.Name
}

func (c *ModelConfig) GetID() int64 {
	return c.ID
}

func (c *ModelConfig) String() string {
	return fmt.Sprintf("%s-%d", FormatString(c.Name), c.ID)
}
