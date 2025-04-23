package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api-client"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/fatih/color"
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
		config, err := common.LoadConfig[common.Config](common.GetConfigPath())
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
			fmt.Println("\033[33mUnable to authenticate with this config\033[0m")
			initConfig(cmd, args)
		}
		agentName := ""
		if len(args) > 0 {
			agentName = args[0]
			args = args[1:]
		}
		var agent common.AgentClient
		for _, agentTmp := range config.Agents {
			if agentTmp.Name == agentName {
				agent = agentTmp
				break
			}
		}
		if agent.Name == "" {
			if len(config.Agents) == 0 {
				fmt.Printf("No agents found in the config.\n ")
				return
			}
			fmt.Println("Please provide an agent from this list as the first argument:")
			for _, agentTmp := range config.Agents {
				fmt.Printf("%s (%d)\n", agentTmp.Name, agentTmp.ID)
			}
			return
		}
		instanceName := "Default Instance"
		// @TODO: add a flag for specifying the instance name
		if len(agent.Instances) == 0 {
			fmt.Println("No instances found for that agent")
			return
		}
		instance := agent.Instances[0]
		for _, instanceTmp := range agent.Instances {
			if instanceTmp.Title == instanceName {
				instance = instanceTmp
				break
			}
		}
		// fmt.Printf("Agent name (%+v): %s\n", agent.ID, agent.Name)
		// fmt.Printf("Instance name (%+v): %+v\n\n", instance.ID, instance.Title)
		if len(args) > 0 {
			sendChatRequest(client, instance.ID, args[0])
		}
		for {
			fmt.Fprint(os.Stderr, "\033[5m>\033[0m ") // Blinking '>' written to stderr
			input, _ := reader.ReadString('\n')
			sendChatRequest(client, instance.ID, input)
		}
	},
}

func sendChatRequest(client *api.APIClient, instanceID string, input string) {
	// Start the animation in a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	stopAnimation := make(chan bool)
	go func() {
		defer wg.Done()
		animateLoading(stopAnimation)
	}()

	req := ChatRequest{
		Input: strings.TrimSpace(input),
	}
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Println("\033[33mError marshaling chat request to JSON:\033[0m", err)
		return
	}
	var jsonBuffer bytes.Buffer
	jsonBuffer.Write(jsonBytes)
	res, err := client.CallAPI("GET", "/instances/"+instanceID+"/chat", &jsonBuffer)
	if err != nil {
		fmt.Printf("\033[33m\nError sending chat request: %v\n\033[0m", err)
		// Stop the animation
		stopAnimation <- true
		wg.Wait()
		return
	}
	defer res.Close()

	// Stop the animation
	stopAnimation <- true
	wg.Wait()

	body, err := io.ReadAll(res)
	if err != nil {
		fmt.Printf("\033[33mError reading response body while sending chat request: %v\n\033[0m", err)
		return
	}
	var response ChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("\033[33mError unmarshaling response body while sending chat request: \n\t%v\n\033[0m%s\n", err, string(body))
		return
	}
	fmt.Fprint(os.Stderr, color.GreenString(">"))
	fmt.Printf("%s\n", color.GreenString(response.Answer))
}

func animateLoading(stop chan bool) {
	frames := []string{"...", "..", "."}
	for {
		for _, frame := range frames {
			select {
			case <-stop:
				// Clear the line and exit the animation
				fmt.Print("\r\033[K")
				return
			default:
				// Print the current frame
				fmt.Print("\r\033[K")
				fmt.Printf("\r%s", frame)
				time.Sleep(500 * time.Millisecond) // Adjust speed as needed
			}
		}
	}
}

type ChatRequest struct {
	Input string `json:"input"`
}

type ChatResponse struct {
	Answer string `json:"answer"`
}
