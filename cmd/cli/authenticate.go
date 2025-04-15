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
	"github.com/howeyc/gopass"
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
		config, err := common.LoadConfig[common.Config](common.GetConfigPath())
		if err != nil || config == nil || config.GagarinURL == "" {
			fmt.Printf("Enter Gagarin URL (default: %s): ", defaultGagarinURL)
			input, _ := reader.ReadString('\n')
			gagarinURL = strings.TrimSpace(input)
			if gagarinURL == "" {
				gagarinURL = defaultGagarinURL
			}
			if config == nil {
				config = &common.Config{}
			}
			config.GagarinURL = gagarinURL
			fmt.Print("Do you have an account? (y/n): ")
			input, _ = reader.ReadString('\n')
			if strings.TrimSpace(input) == "n" {
				req := signUpRequest{}
				fmt.Print("Enter invite code (default is empty): ")
				input, _ = reader.ReadString('\n')
				inviteCode := strings.TrimSpace(input)
				if inviteCode != "" {
					req.InviteCode = inviteCode
				}
				req.InviteCode = strings.TrimSpace(input)
				fmt.Print("Enter email: ")
				input, _ = reader.ReadString('\n')
				req.Email = strings.TrimSpace(input)
				fmt.Print("Enter password: ")
				passBytes, err := gopass.GetPasswdMasked()
				if err != nil {
					fmt.Println("Error reading password:", err)
					return
				}
				req.Password = strings.TrimSpace(string(passBytes))
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
				_, err = io.ReadAll(res)
				if err != nil {
					fmt.Printf("Error reading response body while signing up: %v\n", err)
					return
				}
				fmt.Printf("\033[0;32mSigned up successfully!\033[0m\n")
				// fmt.Println(string(body))
				config.Email = req.Email
				config.Password = req.Password
				err = common.SaveConfig[common.Config](config, common.GetConfigPath())
				if err != nil {
					fmt.Println("Error saving config:", err)
					return
				}
			} else {
				fmt.Print("Enter email: ")
				input, _ = reader.ReadString('\n')
				config.Email = strings.TrimSpace(input)
				fmt.Print("Enter password: ")
				passBytes, err := gopass.GetPasswdMasked()
				if err != nil {
					fmt.Println("Error reading password:", err)
					return
				}
				config.Password = strings.TrimSpace(string(passBytes))
				client := api.NewAPIClient(config, gagarinURL)
				err = client.Authenticate()
				if err != nil {
					fmt.Printf("\033[0;31mError authenticating:\033[0m +%v\n", err.Error())
					return
				}
				err = common.SaveConfig[common.Config](config, common.GetConfigPath())
				if err != nil {
					fmt.Printf("\033[0;31mError saving config:\033[0m +%v\n", err.Error())
					return
				}
			}
		}
		client := api.NewAPIClient(config, gagarinURL)
		err = client.Authenticate()
		if err != nil {
			fmt.Println(err)
			fmt.Printf("\033[0;31mUnable to authenticate with this config\033[0m\n")
			initConfig(cmd, args)
			return
		}
		fmt.Printf("\033[0;32mAuthenticated successfully!\033[0m\n")
	},
}
