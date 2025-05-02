package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"

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
	// Disable persistent pre-run
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
		initConfig(cmd, args)
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
	Use:   "validate file [file1 file2 ...]",
	Short: "Validate a Fyodorov configuration",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml", "yml"}, cobra.ShellCompDirectiveFilterFileExt
	},
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		for _, arg := range args {
			wg.Add(1)
			go func(arg string) {
				defer wg.Done()
				validateYamlFile(arg)
			}(arg)
		}
		wg.Wait()
	},
}

func validateYamlFile(filepath string) {
	// Load the config from the file
	config, err := common.LoadConfig[common.FyodorovConfig](filepath)
	if err != nil {
		fmt.Printf("\033[33mError loading fyodorov config '%s' from file: %v\n\033[0m", filepath, err)
		return
	}

	// Load the file directly
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("\033[33mError opening fyodorov config file '%s': %v\n\033[0m", filepath, err)
		return
	}

	// Verify there are no other fields in the file
	var cfg common.FyodorovConfig
	dec := yaml.NewDecoder(bytes.NewReader(fileBytes))
	dec.KnownFields(true) // ← reject any unknown fields
	if err := dec.Decode(&cfg); err != nil {
		fmt.Printf("invalid config %s: %v", filepath, err)
		return
	}
	// Validate the config
	if err := config.Validate(); err != nil {
		fmt.Printf("Fyodorov config '%s' is invalid: %v\n", filepath, err)
		return
	}

	fmt.Printf("\033[36mFyodorov config '%s' is valid\n\033[0m", filepath)
}
