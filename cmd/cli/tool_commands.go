package main

import (
	"fmt"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	cmdTool.AddCommand(createToolCmd)
	cmdTool.AddCommand(validateToolCmd)
	cmdTool.AddCommand(deployToolCmd)
}

// Tool commands
var cmdTool = &cobra.Command{
	Use:   "tool [create|validate|deploy] file",
	Short: "Perform operations on tool configurations",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Perform operations on tool configurations")
	},
}

var createToolCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Create a tool configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]
		// Create empty instance of ToolConfig
		toolConfig := common.ToolConfig{}
		toolConfig.SaveToFile(fileName)
	},
}

var validateToolCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a tool configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := common.LoadToolConfig(args[0])
		if err != nil {
			fmt.Printf("Error loading tool config from file %s: %v\n", args[0], err)
			return
		}
		fmt.Println("Tool config is valid")
	},
}

var deployToolCmd = &cobra.Command{
	Use:   "deploy [file]",
	Short: "Deploy a tool configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// use global config to create an api client, authenticate, and deploy the tool
		client := api.NewAPIClient(config, config.TsiolkovskyURL)
		err := client.Authenticate()
		if err != nil {
			fmt.Println("Error authenticating:", err)
			return
		}
		toolConfig, err := common.LoadToolConfig(args[0])
		if err != nil {
			fmt.Printf("Error loading tool config from file %s: %v\n", args[0], err)
			return
		}
		// marshall tool config to yaml
		yamlBytes, err := yaml.Marshal(toolConfig)
		if err != nil {
			fmt.Printf("Error marshaling tool config to yaml: %v\n", err)
			return
		}
		// deploy tool
		_, err = client.CallAPI("POST", "/tools/yaml", yamlBytes)
		if err != nil {
			fmt.Printf("Error deploying tool: %v\n", err)
			return
		}
		fmt.Println("Tool deployed successfully")
	},
}
