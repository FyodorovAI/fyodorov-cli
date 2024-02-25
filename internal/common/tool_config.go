package common // Adjust the package name based on your actual package structure

import "fmt" // ToolConfig represents the configuration structure for a tool
type ToolConfig struct {
	NameForHuman        string `json:"name_for_human" yaml:"name_for_human"`
	NameForAI           string `json:"name_for_ai" yaml:"name_for_model"` // Adjust based on your structure
	DescriptionForHuman string `json:"description_for_human" yaml:"description_for_human"`
	DescriptionForAI    string `json:"description_for_ai" yaml:"description_for_model"` // Adjust based on your structure
	APIType             string `json:"api_type" yaml:"api.type"`
	APIURL              string `json:"api_url" yaml:"api.url"`
	LogoURL             string `json:"logo_url" yaml:"logo_url"`
	ContactEmail        string `json:"contact_email" yaml:"contact_email"`
	LegalInfoURL        string `json:"legal_info_url" yaml:"legal_info_url"`
	// Include fields from the 'auth' structure if needed
	Auth struct {
		AuthorizationType string `json:"authorization_type" yaml:"auth.authorization_type"`
		Type              string `json:"type" yaml:"auth.type"`
	} `json:"auth" yaml:"auth"`
}

// LoadToolConfig loads tool configuration from a file.
// The format is determined by the file extension: .json or .yaml
func LoadToolConfig(filename string) (*ToolConfig, error) {
	config, err := LoadConfig[ToolConfig](filename)
	if err != nil {
		fmt.Printf("Error loading tool config from file %s: %v\n", filename, err)
		return nil, err
	}

	return config, nil
}

// The format is determined by the file extension: .json or .yaml
func (config *ToolConfig) SaveToFile(filename string) error {
	SaveConfig[ToolConfig](config, filename)

	return nil
}
