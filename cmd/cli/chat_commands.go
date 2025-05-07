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

	"golang.org/x/exp/slices"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api-client"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	// Add command-specific flags
	chatCmd.Flags().String("agent", "", "Specify the agent name")
	chatCmd.Flags().String("instance", "", "Specify the instance name")
	rootCmd.AddCommand(chatCmd)
}

// Fyodorov commands
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with an instance of an agent",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		agents := GetResources().Agents
		if len(args) == 0 {
			agentNames := make([]string, len(agents))
			for i, agent := range agents {
				agentNames[i] = agent.Name
			}
			return agentNames, cobra.ShellCompDirectiveNoFileComp
		} else if len(args) == 1 {
			agentName := args[0]
			var agent common.Agent
			for _, agentTmp := range agents {
				if agentTmp.Name == agentName {
					agent = agentTmp
					break
				}
			}
			instances := GetResources().Instances
			agentInstances := GetAgentInstances(instances, agent.ID)
			instanceNames := make([]string, len(instances))
			for i, instance := range agentInstances {
				instanceNames[i] = instance.Name
			}
			return instanceNames, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		if !v.IsSet("gagarin-url") {
			fmt.Println("Enter Gagarin URL:")
			input, _ := reader.ReadString('\n')
			v.Set("gagarin-url", strings.TrimSpace(input))
		}
		client, err := api.NewAPIClient(v, "")
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
		var agent common.Agent
		agents := GetResources().Agents
		instances := GetResources().Instances
		for _, agentTmp := range agents {
			if agentTmp.Name == agentName {
				agent = agentTmp
				break
			}
		}
		if agent.Name == "" {
			fmt.Println("Please provide an agent from this list as the first argument:")
			for _, agentTmp := range agents {
				fmt.Printf("%s (%d)\n", agentTmp.Name, agentTmp.ID)
			}
			return
		}
		instances = GetResources().Instances
		instances = slices.DeleteFunc(instances, func(instance common.Instance) bool {
			return instance.AgentId != agent.ID
		})
		// @TODO: add a flag for specifying the instance name
		if len(instances) == 0 {
			fmt.Println("No instances found for that agent")
			client, err := api.NewAPIClient(v, "")
			if err != nil {
				fmt.Printf("\033[33mError creating API client:\033[0m +%v\n", err.Error())
				return
			}
			err = client.Authenticate()
			if err != nil {
				fmt.Printf("\033[33mError authenticating:\033[0m %v\n", err.Error())
				return
			}
			instance := common.Instance{
				AgentId: agent.ID,
				Title:   fmt.Sprintf("Default Instance (%d)", agent.ID),
			}
			jsonBytes, err := json.Marshal(instance)
			if err != nil {
				fmt.Printf("\033[33mError marshaling instance to JSON:\033[0m %v\n", err.Error())
				return
			}
			res, err := client.CallAPI("POST", "/instances", bytes.NewBuffer(jsonBytes))
			if err != nil {
				fmt.Printf("\033[33mError creating instance:\033[0m %v\n", err.Error())
				return
			}
			defer res.Close()
			fmt.Printf("\033[36mCreated instance (%s)\033[0m\n", instance.Title)
			err = json.NewDecoder(res).Decode(&instance)
			if err != nil {
				fmt.Printf("\033[33mError decoding response:\033[0m %v\n", err.Error())
				return
			}
			instances = append(instances, instance)
		}
		instance := instances[0]
		instanceName := instance.Title
		if len(args) > 1 {
			instanceName = args[1]
		}
		for _, instanceTmp := range instances {
			if instanceTmp.Title == instanceName {
				instance = instanceTmp
				break
			}
		}
		// fmt.Printf("Agent name (%+v): %s\n", agent.ID, agent.Name)
		// fmt.Printf("Instance name (%+v): %+v\n\n", instance.ID, instance.Title)
		if len(args) > 0 {
			err = sendChatRequest(client, instance.ID, args[0])
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		}
		for {
			fmt.Fprint(os.Stderr, "\033[5m>\033[0m ") // Blinking '>' written to stderr
			input, _ := reader.ReadString('\n')
			err = sendChatRequest(client, instance.ID, input)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		}
	},
}

func GetAgentInstances(instances []common.Instance, agentID int64) []common.InstanceClient {
	instanceClients := make([]common.InstanceClient, 0)
	for _, instance := range instances {
		if instance.AgentId == agentID {
			instanceClient := common.InstanceClient{
				ID:   instance.ID,
				Name: instance.Title,
			}
			instanceClients = append(instanceClients, instanceClient)
		}
	}
	return instanceClients
}

func sendChatRequest(client *api.APIClient, instanceID int64, input string) error {
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
		return err
	}
	var jsonBuffer bytes.Buffer
	jsonBuffer.Write(jsonBytes)
	instanceIDStr := fmt.Sprintf("%d", instanceID)
	res, err := client.CallAPI("GET", "/instances/"+instanceIDStr+"/chat", &jsonBuffer)
	if err != nil {
		fmt.Printf("\033[33m\nError sending chat request: %v\n\033[0m", err)
		// Stop the animation
		stopAnimation <- true
		wg.Wait()
		return err
	}
	defer res.Close()

	// Stop the animation
	stopAnimation <- true
	wg.Wait()

	body, err := io.ReadAll(res)
	if err != nil {
		fmt.Printf("\033[33mError reading response body while sending chat request: %v\n\033[0m", err)
		return err
	}
	var response ChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("\033[33mError unmarshaling response body while sending chat request: \n\t%v\n\033[0m%s\n", err, string(body))
		return err
	}
	fmt.Fprint(os.Stderr, color.GreenString(">"))
	fmt.Printf("%s\n", color.GreenString(response.Answer))
	return nil
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
