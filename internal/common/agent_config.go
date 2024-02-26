package common

type AgentConfig struct {
	ProviderID          string   `json:"provider_id" yaml:"provider_id"`
	Model               string   `json:"model" yaml:"model"`
	NameForHuman        string   `json:"name" yaml:"name"`
	DescriptionForHuman string   `json:"description" yaml:"description"`
	Prompt              string   `json:"prompt" yaml:"prompt"`
	PromptSize          int      `json:"prompt_size" yaml:"prompt_size"`
	Tools               []string `json:"tools" yaml:"tools"`
	Rag                 []string `json:"rag" yaml:"rag"`
}
