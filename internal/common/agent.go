package common

import "fmt"

var (
	MAX_LENGTH_NAME        = 80
	MAX_LENGTH_DESCRIPTION = 280
)

type Agent struct {
	ID          int64    `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
	Model       string   `json:"model" yaml:"model,omitempty"`
	Description string   `json:"description" yaml:"description,omitempty"`
	Prompt      string   `json:"prompt" yaml:"prompt,omitempty"`
	Tools       []string `json:"tools" yaml:"tools,omitempty"`
	Rag         []string `json:"rag" yaml:"rag,omitempty"`
}

func (c *Agent) Validate() error {
	if c.Model == "" {
		return fmt.Errorf("model is required")
	}
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(c.Name) > MAX_LENGTH_NAME {
		return fmt.Errorf("name cannot exceed %d characters", MAX_LENGTH_NAME)
	}
	if len(c.Description) > MAX_LENGTH_DESCRIPTION {
		return fmt.Errorf("description cannot exceed %d characters", MAX_LENGTH_DESCRIPTION)
	}
	if c.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	return nil
}

func (c *Agent) GetID() int64 {
	return c.ID
}

func (c *Agent) String() string {
	return fmt.Sprintf("%s-%d", FormatString(c.Name), c.ID)
}
