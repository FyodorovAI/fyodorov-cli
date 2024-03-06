package common

// Tool represents the configuration structure for a tool
type Tool struct {
	NameForHuman        string `json:"name" yaml:"name,omitempty"`
	NameForAI           string `json:"name_for_ai" yaml:"name_for_ai,omitempty"`
	DescriptionForHuman string `json:"description" yaml:"description,omitempty"`
	DescriptionForAI    string `json:"description_for_ai" yaml:"description_for_ai,omitempty"`
	API                 struct {
		Type string `json:"type" yaml:"type,omitempty"`
		URL  string `json:"url" yaml:"url,omitempty"`
	} `json:"api" yaml:"api,omitempty"`
	LogoURL      string `json:"logo_url" yaml:"logo_url,omitempty"`
	ContactEmail string `json:"contact_email" yaml:"contact_email,omitempty"`
	LegalInfoURL string `json:"legal_info_url" yaml:"legal_info_url,omitempty"`
	// Include fields from the 'auth' structure if needed
	Auth struct {
		Type string `json:"type" yaml:"type,omitempty"`
	} `json:"auth" yaml:"auth,omitempty"`
}
