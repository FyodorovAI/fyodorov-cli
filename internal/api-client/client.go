// client.go
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	httpClient = &http.Client{
		Timeout: 20 * time.Second,
	}
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
	if config.TsiolkovskyURL != "" {
		host = config.TsiolkovskyURL
	}
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
	url := c.BaseURL + endpoint
	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, body)
	}
	if err != nil {
		return nil, err
	}
	// Set the necessary headers, for example, Authorization headers
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}
	req.Header.Set("User-Agent", "fyodorov-cli-tool")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		// Handle HTTP errors here
		return nil, fmt.Errorf("[%s] API request error: %s", url, resp.Status)
	}
	return resp.Body, nil
}

func (c *APIClient) GetResources(resourceType *string, v *viper.Viper) (*common.FyodorovConfig, error) {
	config := common.CreateFyodorovConfig(v)
	var response io.ReadCloser
	var err error
	if resourceType == nil {
		response, err = c.CallAPI("GET", "/yaml/", nil)
	} else {
		response, err = c.CallAPI("GET", fmt.Sprintf("/yaml/%s", *resourceType), nil)
	}
	if err != nil {
		return config, err
	}
	defer response.Close()
	body, err := io.ReadAll(response)
	if err != nil {
		return config, err
	}
	dec := yaml.NewDecoder(bytes.NewReader(body))
	// dec.KnownFields(true) // ← reject any unknown fields
	if err := dec.Decode(&config); err != nil {
		fmt.Printf("invalid config: %v", err)
		return config, err
	}
	return config, nil
}
