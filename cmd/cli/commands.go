package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	// rootCmd.AddCommand(copilotCmd)
	rootCmd.AddCommand(validateTemplateCmd)
	rootCmd.AddCommand(configCmd)
}

// Fyodorov commands
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Fyodorov configuration",
	Run: func(cmd *cobra.Command, args []string) {
		configRun = true
		initConfig(cmd, args)
		configRun = false
	},
}

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
		var cfg common.FyodorovConfig
		dec := yaml.NewDecoder(bytes.NewReader(fileBytes))
		dec.KnownFields(true) // ← reject any unknown fields
		if err := dec.Decode(&cfg); err != nil {
			fmt.Printf("invalid config: %v", err)
			return
		}
		// Validate the config
		if err := config.Validate(); err != nil {
			fmt.Printf("Fyodorov config is invalid: %v\n", err)
			return
		}

		fmt.Println("Fyodorov config is valid")
	},
}
