package common

// Define a struct to hold the configuration
type Config struct {
	GagarinURL     string `json:"gagarin_url"`
	TsiolkovskyURL string `json:"tsiolkovsky_url"`
	DostoyevskyURL string `json:"dostoyevsky_url"`
	Email          string `json:"email"`
	Password       string `json:"password"`
}
