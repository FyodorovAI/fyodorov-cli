package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api-client"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "fyodorov [validate|deploy|chat|config|list|remove] file",
		Short: "Fyodorov CLI tool",
	}

	defaultGagarinURL = "https://gagarin.danielransom.com"
	NoCache           bool
	// Init viper instance
	v = common.InitViper()
)

func main() {
	// Define flag variables
	var (
		gagarinURL     string
		tsiolkovskyURL string
		// dostoyevskyURL string
		email    string
		password string
	)

	// Define global flags
	rootCmd.PersistentFlags().StringVarP(&gagarinURL, "gagarin-url", "b", "", "base URL for 'Gagarin'")
	rootCmd.PersistentFlags().StringVarP(&tsiolkovskyURL, "tsiolkovsky-url", "t", "", "base URL for 'Tsiolkovsky'")
	// rootCmd.PersistentFlags().StringVarP(&dostoyevskyURL, "dostoyevsky-url", "t", "", "base URL for 'Dostoyevsky'")
	rootCmd.PersistentFlags().StringVarP(&email, "email", "u", "", "email for authentication")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password for authentication")
	rootCmd.PersistentFlags().BoolVarP(&NoCache, "no-cache", "n", false, "disable cache")

	// Bind flags to Viper
	v.BindPFlag("gagarin-url", rootCmd.PersistentFlags().Lookup("gagarin-url"))
	v.BindPFlag("tsiolkovsky-url", rootCmd.PersistentFlags().Lookup("tsiolkovsky-url"))
	// v.BindPFlag("dostoyevsky-url", rootCmd.PersistentFlags().Lookup("dostoyevsky-url"))
	v.BindPFlag("email", rootCmd.PersistentFlags().Lookup("email"))
	v.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	v.BindPFlag("no-cache", rootCmd.PersistentFlags().Lookup("no-cache"))

	rootCmd.PersistentPreRun = initConfig

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig(cmd *cobra.Command, args []string) {
	fmt.Println("CLI Version:", version)
	configRun := cmd.Use == "config"
	reader := bufio.NewReader(os.Stdin)

	// Prompt for missing values
	if !v.IsSet("gagarin-url") {
		fmt.Printf("Enter Gagarin URL (default: %s): ", defaultGagarinURL)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			v.Set("gagarin-url", input)
		} else {
			v.Set("gagarin-url", defaultGagarinURL)
		}
	}
	if !v.IsSet("tsiolkovsky-url") {
		defaultTsiolkovskyURL := strings.Replace(v.GetString("gagarin-url"), "gagarin", "tsiolkovsky", -1)
		fmt.Printf("Enter Tsiolkovsky URL (default: %s): ", defaultTsiolkovskyURL)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			v.Set("tsiolkovsky-url", input)
		} else {
			v.Set("tsiolkovsky-url", defaultTsiolkovskyURL)
		}
	}
	// if !v.IsSet("dostoyevsky-url") {
	// defaultDostoyevskyURL := strings.Replace(v.GetString("gagarin-url"), "gagarin", "dostoyevsky", -1)
	// fmt.Printf("Enter Dostoyevsky URL (default: %s): ", defaultDostoyevskyURL)
	// input, _ := reader.ReadString('\n')
	// input = strings.TrimSpace(input)
	// if input != "" {
	// 	v.Set("dostoyevsky-url", input)
	// }
	// else {
	// 	v.Set("dostoyevsky-url", defaultDostoyevskyURL)
	// }
	// }
	if !v.IsSet("email") {
		fmt.Print("Enter Email: ")
		input, _ := reader.ReadString('\n')
		v.Set("email", strings.TrimSpace(input))
	}
	if !v.IsSet("password") {
		fmt.Print("Enter Password: ")
		passBytes, err := gopass.GetPasswdMasked()
		if err != nil {
			fmt.Println("Error getting password:", err)
			return
		}
		v.Set("password", strings.TrimSpace(string(passBytes)))
	}

	// Print config
	if configRun {
		fmt.Printf("--Config-------------------------------------\n")
		fmt.Printf("Gagarin URL: %s\n", v.GetString("gagarin-url"))
		fmt.Printf("Tsiolkovsky URL: %s\n", v.GetString("tsiolkovsky-url"))
		// fmt.Printf("Dostoyevsky URL: %s\n", config.DostoyevskyURL)
		fmt.Printf("Email: %s\n", v.GetString("email"))
		// Replace all but first and last letter with '*'
		password := v.GetString("password")
		fmt.Printf(
			"Password: %s\n\n",
			strings.ReplaceAll(
				password,
				password[1:len(password)-2],
				strings.Repeat("*", len(password)-3),
			),
		)
	}

	// Initialize API client
	client, err := api.NewAPIClient(v, "")
	if err != nil {
		fmt.Printf("\033[0;31mError creating API client:\033[0m %v\n", err)
		return
	}

	// Authenticate
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
		fmt.Printf("\033[0;36mAuthenticated config successfully!\033[0m\n")
	}

	// Save the updated configuration
	configPath := common.GetConfigPath()
	if err := v.WriteConfigAs(configPath); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
	}
}
