package common

import (
	"fmt"
	"net/url"
)

type Provider struct {
	ID     int64  `json:"id,omitempty" yaml:"id,omitempty"`
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	URL    string `json:"api_url" yaml:"api_url,omitempty"`
	APIKey string `json:"api_key" yaml:"api_key,omitempty"`
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
	if c.URL != "" {
		_, err := url.Parse(c.URL)
		if err != nil {
			return err
		}
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

func (c *Provider) GetID() int64 {
	return c.ID
}

func (c *Provider) String() string {
	return fmt.Sprintf("%s-%d", FormatString(c.Name), c.ID)
}
