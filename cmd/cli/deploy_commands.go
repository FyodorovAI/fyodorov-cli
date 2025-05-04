package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api-client"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	dryRun bool
	values []string
)

func init() {
	deployTemplateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run")
	deployTemplateCmd.Flags().StringSliceVar(&values, "set", []string{}, "List of key=value pairs (e.g. --set key1=value1,key2=value2)")
	rootCmd.AddCommand(deployTemplateCmd)
}

// Fyodorov commands
var deployTemplateCmd = &cobra.Command{
	Use:   "deploy file [file1 file2 ...]",
	Short: "Deploy a Fyodorov configuration",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml", "yml"}, cobra.ShellCompDirectiveFilterFileExt
	},
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		// Allow deploying multiple configs passed as arguments
		for _, arg := range args {
			wg.Add(1)
			go func(arg string) {
				defer wg.Done()
				deployYamlFile(arg)
			}(arg)
		}
		cache.Update(true)
		wg.Wait()
	},
}

func deployYamlFile(filepath string) {
	FyodorovConfig, err := common.LoadConfig[common.FyodorovConfig](filepath)
	if err != nil {
		fmt.Printf("\033[33mError loading fyodorov config from file %s: %v\n\033[0m", filepath, err)
		return
	}
	// load fyodorov config from values
	if len(values) > 0 {
		FyodorovConfig.ParseKeyValuePairs(values)
	}
	// validate fyodorov config
	err = FyodorovConfig.Validate()
	if err != nil {
		fmt.Printf("\033[33mError validating fyodorov config (%s): %v\n\033[0m", filepath, err)
		return
	}
	// print fyodorov config to stdout
	if dryRun {
		bytes, err := yaml.Marshal(FyodorovConfig)
		if err != nil {
			fmt.Printf("\033[33mError marshaling fyodorov config to yaml: %v\n\033[0m", err)
			return
		}
		// Print the YAML to stdout
		fmt.Printf("\033[36mValidated config %s\033[0m\n", filepath)
		fmt.Printf("---Fyodorov config---\n%s\n", string(bytes))
		return
	}
	// deploy config to Gagarin
	if !dryRun {
		yamlBytes, err := yaml.Marshal(FyodorovConfig)
		if err != nil {
			fmt.Printf("\033[34mError marshaling config to yaml: %v\n\033[0m", err)
			return
		}
		client, err := api.NewAPIClient(v, "")
		if err != nil {
			return
		}
		err = client.Authenticate()
		if err != nil {
			fmt.Println("\033[33mError authenticating during deploy:\033[0m", err)
			fmt.Println("\033[33mUnable to authenticate with this config\033[0m")
			return
		}
		var yamlBuffer bytes.Buffer
		yamlBuffer.Write(yamlBytes)
		res, err := client.CallAPI("POST", "/yaml", &yamlBuffer)
		if err != nil {
			fmt.Printf("\033[33mError deploying config (%s): %v\n\033[0m", filepath, err.Error())
			return
		}
		defer res.Close()
		body, err := io.ReadAll(res)
		if err != nil {
			fmt.Printf("\033[33mError reading response body while deploying config: %v\n\033[0m", err)
			return
		}
		var bodyResponse BodyResponse
		err = json.Unmarshal(body, &bodyResponse)
		if err != nil {
			fmt.Printf("\033[33mError unmarshaling response body while deploying config (%s): %s\n\t%v\n\033[0m", filepath, string(body), err)
			return
		}
		// Print deployed config
		fmt.Printf("\033[36mDeployed config %s\033[0m\n", filepath)
	}
}

type DatabaseMetadata struct {
	CreatedAt string `json:"created_at"`
	UserID    string `json:"user_id"`
}

type ProviderResponse struct {
	common.Provider
	DatabaseMetadata
	ProviderID int `json:"provider"`
}

type ModelResponse struct {
	common.ModelInfo
	DatabaseMetadata
	common.ModelConfig
	ProviderID int `json:"provider"`
}

type AgentResponse struct {
	common.Agent
	DatabaseMetadata
}

type InstanceResponse struct {
	common.Instance
	DatabaseMetadata
	AgentID int64 `json:"agent_id"`
}

type ToolResponse struct {
	common.Tool
	DatabaseMetadata
}

type BodyResponse struct {
	Providers []ProviderResponse `json:"providers"`
	Models    []ModelResponse    `json:"models"`
	Agents    []AgentResponse    `json:"agents"`
	Tools     []ToolResponse     `json:"tools"`
	Instances []InstanceResponse `json:"instances"`
}
