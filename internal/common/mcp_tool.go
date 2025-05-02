package common

import (
	"fmt"
	"net/url"
	"time"
)

// MCPTool mirrors the 'mcp_tools' table schema.
type MCPTool struct {
	ID           int64                  `json:"id,omitempty" yaml:"id,omitempty"`
	Name         string                 `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt    time.Time              `json:"created_at,omitempty" yaml:"-" db:"created_at,omitempty"` // timestamptz
	UpdatedAt    time.Time              `json:"updated_at,omitempty" yaml:"-" db:"updated_at,omitempty"` // timestamptz
	Handle       string                 `json:"handle,omitempty" yaml:"handle,omitempty"`                // text
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`      // text
	LogoURL      string                 `json:"logo_url,omitempty" yaml:"logo_url,omitempty"`            // text
	UserID       string                 `json:"user_id,omitempty" yaml:"-"`                              // uuid
	Public       bool                   `json:"public,omitempty" yaml:"-"`                               // bool
	APIType      string                 `json:"api_type,omitempty" yaml:"api_type,omitempty"`            // text
	APIURL       string                 `json:"api_url,omitempty" yaml:"api_url,omitempty"`              // text
	AuthMethod   string                 `json:"auth_method,omitempty" yaml:"auth_method,omitempty"`      // text
	AuthInfo     map[string]interface{} `json:"auth_info,omitempty" yaml:"auth_info,omitempty"`          // jsonb
	Capabilities map[string]interface{} `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`    // jsonb
	HealthStatus string                 `json:"health_status,omitempty" yaml:"-"`                        // text
	UsageNotes   string                 `json:"usage_notes,omitempty" yaml:"usage_notes,omitempty"`      // text
}

func (t *MCPTool) Validate() error {
	if t.Handle == "" {
		return fmt.Errorf("tool handle is required")
	}
	if t.LogoURL != "" {
		if _, err := url.Parse(t.LogoURL); err != nil {
			return fmt.Errorf("invalid logo URL: %w", err)
		}
	}
	if t.APIURL != "" {
		if _, err := url.Parse(t.APIURL); err != nil {
			return fmt.Errorf("invalid API URL: %w", err)
		}
	}
	return nil
}

func (t *MCPTool) GetID() int64 {
	return t.ID
}

func (t *MCPTool) String() string {
	var name string
	if t.Handle != "" {
		name = t.Handle
	} else if t.Name != "" {
		name = t.Name
	} else {
		name = "tool"
	}
	return fmt.Sprintf("%s-%d", FormatString(name), t.ID)
}
