package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

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
	Use:   "deploy [file]",
	Short: "Deploy a Fyodorov configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// load fyodorov config from file
		FyodorovConfig, err := common.LoadConfig[common.FyodorovConfig](args[0])
		if err != nil {
			fmt.Printf("Error loading fyodorov config from file %s: %v\n", args[0], err)
			return
		}
		// load fyodorov config from values
		if len(values) > 0 {
			FyodorovConfig.ParseKeyValuePairs(values)
		}
		// validate fyodorov config
		err = FyodorovConfig.Validate()
		if err != nil {
			fmt.Printf("Error validating fyodorov config: %v\n", err)
			return
		}
		// print fyodorov config to stdout
		if dryRun {
			bytes, err := yaml.Marshal(FyodorovConfig)
			if err != nil {
				fmt.Printf("Error marshaling fyodorov config to yaml: %v\n", err)
				return
			}
			// Print the YAML to stdout
			fmt.Println("Validated config")
			fmt.Printf("---Fyodorov config---\n%s\n", string(bytes))
			return
		}
		// deploy config to Gagarin
		if !dryRun {
			yamlBytes, err := yaml.Marshal(FyodorovConfig)
			if err != nil {
				fmt.Printf("Error marshaling config to yaml: %v\n", err)
				return
			}
			client := api.NewAPIClient(config, config.GagarinURL)
			err = client.Authenticate()
			if err != nil {
				fmt.Println("Error authenticating:", err)
				return
			}
			var yamlBuffer bytes.Buffer
			yamlBuffer.Write(yamlBytes)
			res, err := client.CallAPI("POST", "/yaml", &yamlBuffer)
			if err != nil {
				fmt.Printf("Error deploying config: %v\n", err)
				return
			}
			defer res.Close()
			body, err := io.ReadAll(res)
			if err != nil {
				fmt.Printf("Error reading response body while deploying config: %v\n", err)
				return
			}
			var bodyResponse BodyResponse
			err = json.Unmarshal(body, &bodyResponse)
			if err != nil {
				fmt.Printf("Error unmarshaling response body while deploying config: %s\n\t%v\n", string(body), err)
				return
			}
			cliConfig, err := common.LoadConfig[common.Config](common.GetConfigPath())
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}
			for _, agent := range bodyResponse.Agents {
				if checkIfAgentPresent(agent.ID, cliConfig.Agents) {
					continue
				}
				cliConfig.Agents = append(cliConfig.Agents, common.AgentClient{
					ID:        agent.ID,
					Name:      agent.Name,
					Instances: getInstanceForAgent(agent.ID, bodyResponse.Instances),
				})
			}
			common.SaveConfig[common.Config](cliConfig, common.GetConfigPath())
			// Print deployed config
			fmt.Println("Deployed config")
		}
	},
}

func checkIfAgentPresent(agentID int, agents []common.AgentClient) bool {
	for _, agent := range agents {
		if agent.ID == agentID {
			return true
		}
	}
	return false
}

func getInstanceForAgent(agentID int, instances []InstanceResponse) []common.InstanceClient {
	res := make([]common.InstanceClient, 0)
	for _, instance := range instances {
		if instance.AgentID == fmt.Sprint(agentID) {
			res = append(res, common.InstanceClient{
				ID:    instance.ID,
				Title: instance.Title,
			})
		}
	}
	return res
}

type DatabaseMetadata struct {
	ID        int    `json:"id"`
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
	DatabaseMetadata
	Title   string `json:"title"`
	ID      string `json:"id"`
	AgentID string `json:"agent_id"`
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
