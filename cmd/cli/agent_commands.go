package main

import (
	"fmt"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	cmdAgent.AddCommand(createAgentCmd)
	cmdAgent.AddCommand(validateAgentCmd)
	cmdAgent.AddCommand(deployAgentCmd)
}

// Agent commands
var cmdAgent = &cobra.Command{
	Use:   "agent [create|validate|deploy] file",
	Short: "Perform operations on agent configurations",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Implement your handler here based on args
	},
}

var createAgentCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Create an agent configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]
		// Create empty instance of AgentConfig
		agentConfig := common.AgentConfig{}
		agentConfig.SaveToFile(fileName)

		fmt.Printf("Created agent config file %s\n", fileName)
	},
}

var validateAgentCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate an agent configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Load the config from the file and verify there are no other fields in the file
		_, err := common.LoadAgentConfigFromFile(args[0])
		if err != nil {
			fmt.Printf("Error loading agent config from file %s: %v\n", args[0], err)
			return
		}

		fmt.Println("Agent config is valid")
	},
}

var deployAgentCmd = &cobra.Command{
	Use:   "deploy [file]",
	Short: "Deploy an agent configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// use global config to create an api client, authenticate, and deploy the agent
		client := api.NewAPIClient(config, config.GagarinURL)
		err := client.Authenticate()
		if err != nil {
			fmt.Println("Error authenticating:", err)
			return
		}
		// load agent config from file
		agentConfig, err := common.LoadAgentConfigFromFile(args[0])
		if err != nil {
			fmt.Printf("Error loading agent config from file %s: %v\n", args[0], err)
			return
		}
		// marshall agent config to yaml
		yamlBytes, err := yaml.Marshal(agentConfig)
		if err != nil {
			fmt.Printf("Error marshaling agent config to yaml: %v\n", err)
			return
		}
		// deploy agent config
		client.CallAPI("/agents/from-yaml", "POST", yamlBytes)
	},
}
