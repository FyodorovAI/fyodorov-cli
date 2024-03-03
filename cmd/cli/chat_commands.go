package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api-client"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(chatCmd)
}

// Fyodorov commands
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Manage Fyodorov configuration",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		config, err := common.LoadConfig[common.Config](configPath)
		if err != nil || config == nil || config.GagarinURL == "" {
			fmt.Println("Enter Gagarin URL:")
			input, _ := reader.ReadString('\n')
			gagarinURL = strings.TrimSpace(input)
			if config == nil {
				config = &common.Config{}
			}
			config.GagarinURL = gagarinURL
		}
		client := api.NewAPIClient(config, gagarinURL)
		err = client.Authenticate()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Unable to authenticate with this config")
			initConfig(cmd, args)
		}
		agentName := "agent"
		if len(args) > 0 {
			agentName = args[0]
		}
		fmt.Printf("Agent name: %s\n", agentName)
	},
}
