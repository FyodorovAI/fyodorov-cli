package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api-client"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)
}

type signUpRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	InviteCode string `json:"invite_code,omitempty"`
}

// Fyodorov commands
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Fyodorov authentication: sign up, log in, etc.",
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
			fmt.Println("Do you have an account? (y/n)")
			input, _ = reader.ReadString('\n')
			if strings.TrimSpace(input) == "y" {
				req := signUpRequest{}
				fmt.Println("Enter invite code:")
				input, _ = reader.ReadString('\n')
				req.InviteCode = strings.TrimSpace(input)
				fmt.Println("Enter email:")
				input, _ = reader.ReadString('\n')
				req.Email = strings.TrimSpace(input)
				fmt.Println("Enter password:")
				input, _ = reader.ReadString('\n')
				req.Password = strings.TrimSpace(input)
				// marshal request to json
				jsonBytes, err := json.Marshal(req)
				if err != nil {
					fmt.Println("Error marshaling sign up request to JSON:", err)
					return
				}
				var jsonBuffer bytes.Buffer
				jsonBuffer.Write(jsonBytes)
				// call API
				client := api.NewAPIClient(config, gagarinURL)
				res, err := client.CallAPI("POST", "/users/sign_up", &jsonBuffer)
				if err != nil {
					fmt.Printf("Error signing up: %v\n", err)
					return
				}
				defer res.Close()
				body, err := io.ReadAll(res)
				if err != nil {
					fmt.Printf("Error reading response body while signing up: %v\n", err)
					return
				}
				fmt.Println("Signed up with response:")
				fmt.Println(string(body))
				config.Email = req.Email
				config.Password = req.Password
				err = common.SaveConfig[common.Config](config, configPath)
				if err != nil {
					fmt.Println("Error saving config:", err)
					return
				}
			}
		}
		client := api.NewAPIClient(config, gagarinURL)
		err = client.Authenticate()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Unable to authenticate with this config")
			initConfig(cmd, args)
		}
	},
}
