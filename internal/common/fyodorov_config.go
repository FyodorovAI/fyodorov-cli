package common

var DEFAULT_VERSION = "0.0.1"

type FyodorovConfig struct {
	Version   string         `json:"version" yaml:"version"`
	Providers *[]Provider    `json:"providers,omitempty" yaml:"providers,omitempty"`
	Models    *[]ModelConfig `json:"models,omitempty" yaml:"models,omitempty"`
	Agents    *[]Agent       `json:"agents,omitempty" yaml:"agents,omitempty"`
	Tools     *[]Tool        `json:"tools,omitempty" yaml:"tools,omitempty"`
}

func CreateFyodorovConfig() *FyodorovConfig {
	return &FyodorovConfig{
		Version: DEFAULT_VERSION,
		Agents:  nil,
		Tools:   nil,
	}
}
