package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api-client"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var (
	config  *common.Config
	rootCmd = &cobra.Command{
		Use:   "fyodorov [validate|deploy] file",
		Short: "Fyodorov CLI tool",
	}
)

// Define global flags
var (
	gagarinURL string
	// tsiolkovskyURL string
	// dostoyevskyURL string
	email             string
	password          string
	configRun         bool
	defaultGagarinURL = "https://gagarin.danielransom.com"
)

func main() {
	config = &common.Config{}

	// Define global flags
	rootCmd.PersistentFlags().StringVarP(&gagarinURL, "gagarin-url", "b", "", "base URL for 'Gagarin'")
	// rootCmd.PersistentFlags().StringVarP(&tsiolkovskyURL, "tsiolkovsky-url", "t", "", "base URL for 'Tsiolkovsky'")
	// rootCmd.PersistentFlags().StringVarP(&dostoyevskyURL, "dostoyevsky-url", "t", "", "base URL for 'Dostoyevsky'")
	rootCmd.PersistentFlags().StringVarP(&email, "email", "u", "", "email for authentication")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password for authentication")

	rootCmd.PersistentPreRun = initConfig

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig(cmd *cobra.Command, args []string) {
	if cmd.Use == "auth" {
		return
	}
	var err error
	// Load from config if flags are not provided
	// if !configRun && (gagarinURL == "" && tsiolkovskyURL == "" && dostoyevskyURL == "") || email == "" || password == "" {
	if !configRun && gagarinURL == "" || email == "" || password == "" {
		config, err = common.LoadConfig[common.Config](common.GetConfigPath())
		if err != nil {
			fmt.Println("No config file found")
		}

		if gagarinURL == "" && config != nil {
			gagarinURL = config.GagarinURL
		}
		// if tsiolkovskyURL == "" && config != nil {
		// 	tsiolkovskyURL = config.TsiolkovskyURL
		// }
		// if dostoyevskyURL == "" && config != nil {
		// 	dostoyevskyURL = config.DostoyevskyURL
		// }
		if email == "" && config != nil {
			email = config.Email
		}
		if password == "" && config != nil {
			password = config.Password
		}
	}

	// If still missing, prompt the user
	reader := bufio.NewReader(os.Stdin)
	if gagarinURL == "" {
		fmt.Printf("Enter Gagarin URL (default: %s): ", defaultGagarinURL)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		gagarinURL = input
		if gagarinURL == "" {
			gagarinURL = defaultGagarinURL
		}
	}
	// if tsiolkovskyURL == "" {
	// 	defaultTsiolkovskyURL := strings.Replace(gagarinURL, "gagarin", "tsiolkovsky", -1)
	// 	fmt.Printf("Enter Tsiolkovsky URL (default: %s): ", defaultTsiolkovskyURL)
	// 	tsiolkovskyURL, _ = reader.ReadString('\n')
	// 	tsiolkovskyURL = strings.TrimSpace(tsiolkovskyURL)
	// 	if tsiolkovskyURL == "" {
	// 		tsiolkovskyURL = defaultTsiolkovskyURL
	// 	}
	// }
	// if dostoyevskyURL == "" {
	// 	defaultDostoyevskyURL := strings.Replace(tsiolkovskyURL, "gagarin", "dostoyevsky", -1)
	// 	fmt.Printf("Enter Dostoyevsky URL (default: %s): ", defaultDostoyevskyURL)
	// 	dostoyevskyURL, _ = reader.ReadString('\n')
	// 	dostoyevskyURL = strings.TrimSpace(dostoyevskyURL)
	// 	if dostoyevskyURL == "" {
	// 		dostoyevskyURL = defaultDostoyevskyURL
	// 	}
	// }
	if email == "" {
		fmt.Print("Enter Email: ")
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)
	}
	if password == "" {
		fmt.Print("Enter Password: ")
		passBytes, err := gopass.GetPasswdMasked()
		if err != nil {
			fmt.Println("Error getting password:", err)
			return
		}
		password = strings.TrimSpace(string(passBytes))
	}

	if config == nil {
		config = &common.Config{}
	}
	config.GagarinURL = gagarinURL
	// config.TsiolkovskyURL = tsiolkovskyURL
	// config.DostoyevskyURL = dostoyevskyURL
	config.Email = email
	config.Password = password

	// Print config
	if cmd.Use == "config" && configRun {
		fmt.Printf("--Config-------------------------------------\n")
		fmt.Printf("Gagarin URL: %s\n", config.GagarinURL)
		// fmt.Printf("Tsiolkovsky URL: %s\n", config.TsiolkovskyURL)
		// fmt.Printf("Dostoyevsky URL: %s\n", config.DostoyevskyURL)
		fmt.Printf("Email: %s\n", config.Email)
		// Replace all but first and last letter with '*'
		fmt.Printf(
			"Password: %s\n\n",
			strings.ReplaceAll(
				config.Password,
				config.Password[1:len(config.Password)-2],
				strings.Repeat("*", len(config.Password)-3),
			),
		)
	}

	// Initialize API client
	client := api.NewAPIClient(config, "")
	// Authenticate if necessary
	err = client.Authenticate()
	if err != nil && !configRun {
		fmt.Printf("\033[0;31mError authenticating:\033[0m %v\n", err)
		return
	} else if err != nil {
		fmt.Printf("\033[0;31mUnable to authenticate with this config\033[0m\n")
		fmt.Printf("\033[0;31mPlease provide a valid config\033[0m\n")
		fmt.Printf("Invalid config not saved\n")
		return
	} else if configRun {
		fmt.Printf("\033[0;32mAuthenticated successfully!\033[0m\n")
	}

	// Create config directory if it doesn't exist
	dir := filepath.Dir(common.GetConfigPath())
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	// Save the configuration
	err = common.SaveConfig[common.Config](config, common.GetConfigPath())
	if err != nil {
		fmt.Println("Error saving config:", err)
		return
	}
}
