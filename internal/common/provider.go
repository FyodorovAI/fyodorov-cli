package common

import (
	"fmt"
	"net/url"
)

type Provider struct {
	Name   string `json:"name" yaml:"name"`
	URL    string `json:"api_url" yaml:"api_url,omitempty"`
	APIKey string `json:"api_key" yaml:"api_key"`
}

func CreateProvider() *Provider {
	return &Provider{
		Name:   "",
		URL:    "",
		APIKey: "",
	}
}

func (c *Provider) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("provider name is required")
	}
	if c.URL == "" {
		return fmt.Errorf("provider URL is required")
	} else {
		_, err := url.Parse(c.URL)
		if err != nil {
			return err
		}
	}
	if c.APIKey == "" {
		return fmt.Errorf("provider API key is required")
	}
	return nil
}

func Validate(providers *[]Provider) error {
	for _, provider := range *providers {
		err := provider.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
