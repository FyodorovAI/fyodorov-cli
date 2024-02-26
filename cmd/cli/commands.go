package main

import (
	"fmt"
	"io"
	"os"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Fyodorov commands
var copilotCmd = &cobra.Command{
	Use:   "copilot",
	Short: "Ask for help working with Fyodorov",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This feature coming soon!")
	},
}

var validateTemplateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a Fyodorov configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Load the config from the file
		config, err := common.LoadConfig[common.FyodorovConfig](args[0])
		if err != nil {
			fmt.Printf("Error loading fyodorov config from file %s: %v\n", args[0], err)
			return
		}

		// Load the file directly
		fileBytes, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Error opening fyodorov config file %s: %v\n", args[0], err)
			return
		}

		// Verify there are no other fields in the file
		marshalledConfig, err := yaml.Marshal(config)
		if err != nil {
			fmt.Printf("Error marshaling fyodorov config back to yaml: %v\n", err)
			return
		}
		for i := range fileBytes {
			if fileBytes[i] != marshalledConfig[i] {
				fmt.Println("Fyodorov config contains invalid fields")
				return
			}
		}

		fmt.Println("Fyodorov config is valid")
	},
}

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
		// deploy tools
		if FyodorovConfig.Tools != nil && len(*FyodorovConfig.Tools) > 0 {
			yamlBytes, err := yaml.Marshal(FyodorovConfig.Tools)
			if err != nil {
				fmt.Printf("Error marshaling fyodorov tools to yaml: %v\n", err)
				return
			}
			client := api.NewAPIClient(config, config.TsiolkovskyURL)
			err = client.Authenticate()
			if err != nil {
				fmt.Println("Error authenticating:", err)
				return
			}
			res, err := client.CallAPI("/tools/yaml", "POST", yamlBytes)
			if err != nil {
				fmt.Printf("Error deploying tools: %v\n", err)
				return
			}
			defer res.Close()
			body, err := io.ReadAll(res)
			if err != nil {
				fmt.Printf("Error reading response body while deploying tools: %v\n", err)
				return
			}
			fmt.Printf("Deployed tools: %s\n", string(body))
		}

		// deploy agents
		if FyodorovConfig.Agents != nil && len(*FyodorovConfig.Agents) > 0 {
			yamlBytes, err := yaml.Marshal(FyodorovConfig.Agents)
			if err != nil {
				fmt.Printf("Error marshaling agents to yaml: %v\n", err)
				return
			}
			client := api.NewAPIClient(config, config.GagarinURL)
			err = client.Authenticate()
			if err != nil {
				fmt.Println("Error authenticating:", err)
				return
			}
			res, err := client.CallAPI("/agents/from-yaml", "POST", yamlBytes)
			if err != nil {
				fmt.Printf("Error deploying agents: %v\n", err)
				return
			}
			defer res.Close()
			body, err := io.ReadAll(res)
			if err != nil {
				fmt.Printf("Error reading response body while deploying agents: %v\n", err)
				return
			}
			fmt.Printf("Deployed agents: %s\n", string(body))
		}
	},
}