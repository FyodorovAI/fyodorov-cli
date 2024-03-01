package common

// Tool represents the configuration structure for a tool
type Tool struct {
	NameForHuman        string `json:"name" yaml:"name"`
	NameForAI           string `json:"name_for_ai" yaml:"name_for_ai"`
	DescriptionForHuman string `json:"description" yaml:"description"`
	DescriptionForAI    string `json:"description_for_ai" yaml:"description_for_ai"`
	API                 struct {
		Type string `json:"type" yaml:"type"`
		URL  string `json:"url" yaml:"url"`
	} `json:"api" yaml:"api"`
	LogoURL      string `json:"logo_url" yaml:"logo_url"`
	ContactEmail string `json:"contact_email" yaml:"contact_email"`
	LegalInfoURL string `json:"legal_info_url" yaml:"legal_info_url"`
	// Include fields from the 'auth' structure if needed
	Auth struct {
		Type string `json:"type" yaml:"type"`
	} `json:"auth" yaml:"auth"`
}
