package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api-client"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
)

func init() {
	// UpdateOllama()
}

var (
	ollamaPort = 11434
)

type OllamaTagDetails struct {
	ParentModel   string   `json:"parent_model,omitempty"`
	Format        string   `json:"format,omitempty"`
	Family        string   `json:"family,omitempty"`
	Families      []string `json:"families,omitempty"`
	ParameterSize string   `json:"parameter_size,omitempty"`
	Quantization  string   `json:"quantization_level,omitempty"`
}

type OllamaTag struct {
	Name       string           `json:"name,omitempty"`
	Model      string           `json:"model,omitempty"`
	ModifiedAt string           `json:"modified_at,omitempty"`
	Size       int              `json:"size,omitempty"`
	Digest     string           `json:"digest,omitempty"`
	Details    OllamaTagDetails `json:"details,omitempty"`
}

type OllamaTags struct {
	Models []OllamaTag `json:"models,omitempty"`
}

func UpdateOllama() {
	fmt.Println("Starting ollama update...")
	// Run UpdateProvider() every x minutes
	x := 5 // Change x to the desired number of minutes
	duration := time.Duration(x) * time.Minute
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ollama()
		default:
			time.Sleep(10 * time.Second)
		}
	}
}

func ollama() {
	fmt.Println("Ollama function called")
	if !localModelsEnabled {
		return
	}
	fmt.Println("Updating ollama models...")
	// Get API status of localhost ollama
	url := fmt.Sprintf("http://localhost:%d", ollamaPort)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error connecting to ollama: %v\n", err)
		return
	}
	// check status code
	if resp.StatusCode != 200 {
		fmt.Printf("Ollama API status: %s\n", resp.Status)
		return
	}
	// print response body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	fmt.Println(string(body))

	url = fmt.Sprintf("http://localhost:%d/api/tags", ollamaPort)
	resp, err = http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	// Marshall into OllamaTags
	var ollamaTags OllamaTags
	err = json.NewDecoder(resp.Body).Decode(&ollamaTags)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Tags: %+v\n", ollamaTags.Models)
	UpdateProvider()
	for i := range ollamaTags.Models {
		UpdateModel(ollamaTags.Models[i].GetModelConfig())
	}
}

func (t *OllamaTag) GetModelConfig() common.ModelConfig {
	model := common.ModelConfig{
		Name: t.Model,
		ModelInfo: &common.ModelInfo{
			BaseModel:          t.Model,
			InputCostPerToken:  new(float64),
			OutputCostPerToken: new(float64),
		},
	}

	return model
}

func UpdateProvider() {
	// Get API status of localhost ollama
	url := fmt.Sprintf("http://localhost:%d/api/tags", ollamaPort)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	// check status code
	if resp.StatusCode != 200 {
		fmt.Printf("Ollama API status: %s\n", resp.Status)
		return
	}
	defer resp.Body.Close()
	provider := common.Provider{
		Name:   "ollama",
		URL:    "http://localhost:11434",
		APIKey: "",
	}
	// marshal provider to *bytes.Buffer
	providerBytes, err := json.Marshal(provider)
	if err != nil {
		fmt.Println(err)
		return
	}
	var providerBuffer bytes.Buffer
	providerBuffer.Write(providerBytes)
	config, err := common.GetConfig(nil, v)
	if err != nil {
		fmt.Println("No config file found")
	}
	client := api.NewAPIClient(config, config.GagarinURL)
	err = client.Authenticate()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Updating provider: %s\n", string(providerBytes))
	res, err := client.CallAPI("POST", "/providers", &providerBuffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Close()
	body, err := io.ReadAll(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Updated provider:", string(body))
}

func UpdateModel(model common.ModelConfig) {
	config, err := common.GetConfig(nil, v)
	if err != nil {
		fmt.Println("No config file found")
	}
	client := api.NewAPIClient(config, config.GagarinURL)
	err = client.Authenticate()
	if err != nil {
		fmt.Println("Error authenticating during model update:", err)
		return
	}
	// marshall model to *bytes.Buffer
	modelBytes, err := json.Marshal(model)
	if err != nil {
		fmt.Println("Error marshalling model update:", err)
		return
	}
	var modelBuffer bytes.Buffer
	modelBuffer.Write(modelBytes)
	res, err := client.CallAPI("POST", "/models", &modelBuffer)
	if err != nil {
		fmt.Println("Error updating model:", err)
		return
	}
	defer res.Close()
	body, err := io.ReadAll(res)
	if err != nil {
		fmt.Println("Error reading model update response:", err)
		return
	}
	fmt.Println("Updated model:", string(body))
}
