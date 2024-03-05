// client.go
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
)

type APIClient struct {
	BaseURL   string
	Email     string
	Password  string
	AuthToken string
}

func NewAPIClient(config *common.Config, baseURL string) *APIClient {
	var host string
	// if config.DostoyevskyURL != "" {
	// 	host = config.DostoyevskyURL
	// }
	// if config.TsiolkovskyURL != "" {
	// 	host = config.TsiolkovskyURL
	// }
	if config.GagarinURL != "" {
		host = config.GagarinURL
	}
	if baseURL != "" {
		host = baseURL
	}
	return &APIClient{
		BaseURL:  host,
		Email:    config.Email,
		Password: config.Password,
	}
}

// Authenticate method for API client
func (c *APIClient) Authenticate() error {
	// Implement authentication with the API to obtain AuthToken
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(map[string]string{"email": c.Email, "password": c.Password})
	responseBody, err := c.CallAPI("POST", "/users/sign_in", body)
	if err != nil {
		return err
	}

	var response struct {
		Message string `json:"message"`
		JWT     string `json:"jwt"`
	}

	err = json.NewDecoder(responseBody).Decode(&response)
	if err != nil {
		return err
	}

	// fmt.Println(response.Message)
	c.AuthToken = response.JWT
	return nil
}

// CallAPI makes a generic API call. This can be expanded based on your needs.
func (c *APIClient) CallAPI(method, endpoint string, body *bytes.Buffer) (io.ReadCloser, error) {
	// Check if first character of endpoint is '/' and if not add it
	if endpoint[0] != '/' {
		endpoint = "/" + endpoint
	}
	req, err := http.NewRequest(method, c.BaseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	// Set the necessary headers, for example, Authorization headers
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		// Handle HTTP errors here
		return nil, fmt.Errorf("API request error: %s", resp.Status)
	}

	return resp.Body, nil
}
