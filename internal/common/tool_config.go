package common // Adjust the package name based on your actual package structure

// ToolConfig represents the configuration structure for a tool
type ToolConfig struct {
	NameForHuman        string `json:"name" yaml:"name"`
	NameForAI           string `json:"name_for_ai" yaml:"name_for_model"` // Adjust based on your structure
	DescriptionForHuman string `json:"description" yaml:"description"`
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
