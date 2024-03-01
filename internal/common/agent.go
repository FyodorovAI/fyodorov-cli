package common

import "fmt"

type Agent struct {
	Model               string   `json:"model" yaml:"model"`
	NameForHuman        string   `json:"name" yaml:"name"`
	DescriptionForHuman string   `json:"description" yaml:"description"`
	Prompt              string   `json:"prompt" yaml:"prompt"`
	Tools               []string `json:"tools" yaml:"tools"`
	Rag                 []string `json:"rag" yaml:"rag"`
}

func (c *Agent) Validate() error {
	if c.Model == "" {
		return fmt.Errorf("model is required")
	}
	if c.NameForHuman == "" {
		return fmt.Errorf("name is required")
	}
	if len(c.NameForHuman) > 80 {
		return fmt.Errorf("name cannot exceed 80 characters")
	}
	if c.DescriptionForHuman == "" {
		return fmt.Errorf("description is required")
	}
	if len(c.DescriptionForHuman) > 280 {
		return fmt.Errorf("description cannot exceed 200 characters")
	}
	if c.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	return nil
}
