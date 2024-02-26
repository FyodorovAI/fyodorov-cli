package common

var DEFAULT_VERSION = "0.0.1"

type FyodorovConfig struct {
	Version string         `json:"version" yaml:"version"`
	Agents  *[]AgentConfig `json:"agents,omitempty" yaml:"agents,omitempty"`
	Tools   *[]ToolConfig  `json:"tools,omitempty" yaml:"tools,omitempty"`
}

func CreateFyodorovConfig() *FyodorovConfig {
	return &FyodorovConfig{
		Version: DEFAULT_VERSION,
		Agents:  nil,
		Tools:   nil,
	}
}
