package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/spf13/cobra"
)

var config *common.Config

func main() {
	config = &common.Config{}

	// Define global flags
	var (
		gagarinURL     string
		tsiolkovskyURL string
		dostoyevskyURL string
		email          string
		password       string
	)

	var rootCmd = &cobra.Command{
		Use:   "fyodorov",
		Short: "Fyodorov CLI tool",
	}

	// Define global flags
	rootCmd.PersistentFlags().StringVarP(&gagarinURL, "gagarin-url", "b", "", "Base URL for 'Gagarin'")
	rootCmd.PersistentFlags().StringVarP(&tsiolkovskyURL, "tsiolkovsky-url", "t", "", "Base URL for 'Tsiolkovsky'")
	rootCmd.PersistentFlags().StringVarP(&dostoyevskyURL, "dostoyevsky-url", "d", "", "Base URL for 'Dostoyevsky'")
	rootCmd.PersistentFlags().StringVarP(&email, "email", "u", "", "Email for authentication")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for authentication")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Load from config if flags are not provided
		configPath := filepath.Join(os.Getenv("HOME"), ".fyodorov", "config.json")
		if (gagarinURL == "" && tsiolkovskyURL == "" && dostoyevskyURL == "") || email == "" || password == "" {
			config.LoadConfig(configPath)

			if gagarinURL == "" {
				gagarinURL = config.GagarinURL
			}
			if tsiolkovskyURL == "" {
				tsiolkovskyURL = config.TsiolkovskyURL
			}
			if dostoyevskyURL == "" {
				dostoyevskyURL = config.DostoyevskyURL
			}
			if email == "" {
				email = config.Email
			}
			if password == "" {
				password = config.Password
			}
		}

		// If still missing, prompt the user
		reader := bufio.NewReader(os.Stdin)
		if gagarinURL == "" {
			defaultGagarinURL := "https://gagarin.danielransom.com"
			fmt.Printf("Enter Gagarin URL (default: %s): ", defaultGagarinURL)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			gagarinURL = input
			if gagarinURL == "" {
				gagarinURL = defaultGagarinURL
			}
		}
		if tsiolkovskyURL == "" {
			defaultTsiolkovskyURL := strings.Replace(gagarinURL, "gagarin", "tsiolkovsky", -1)
			fmt.Printf("Enter Tsiolkovsky URL (default: %s): ", defaultTsiolkovskyURL)
			tsiolkovskyURL, _ = reader.ReadString('\n')
			tsiolkovskyURL = strings.TrimSpace(tsiolkovskyURL)
			if tsiolkovskyURL == "" {
				tsiolkovskyURL = defaultTsiolkovskyURL
			}
		}
		if dostoyevskyURL == "" {
			defaultDostoyevskyURL := strings.Replace(tsiolkovskyURL, "gagarin", "dostoyevsky", -1)
			fmt.Printf("Enter Dostoyevsky URL (default: %s): ", defaultDostoyevskyURL)
			dostoyevskyURL, _ = reader.ReadString('\n')
			dostoyevskyURL = strings.TrimSpace(dostoyevskyURL)
			if dostoyevskyURL == "" {
				dostoyevskyURL = defaultDostoyevskyURL
			}
		}
		if email == "" {
			fmt.Print("Enter Email: ")
			email, _ = reader.ReadString('\n')
			email = strings.TrimSpace(email)
		}
		if password == "" {
			fmt.Print("Enter Password: ")
			password, _ = reader.ReadString('\n')
			password = strings.TrimSpace(password)
		}

		config.GagarinURL = gagarinURL
		config.TsiolkovskyURL = tsiolkovskyURL
		config.DostoyevskyURL = dostoyevskyURL
		config.Email = email
		config.Password = password

		// Initialize API client
		client := api.NewAPIClient(config, "")

		// Authenticate if necessary
		err := client.Authenticate()
		if err != nil {
			fmt.Println("Error authenticating:", err)
			return
		}

		// Create config directory if it doesn't exist
		dir := filepath.Dir(configPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}

		// Save the configuration
		config.SaveConfig(configPath)
	}

	rootCmd.AddCommand(cmdAgent, cmdTool)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
