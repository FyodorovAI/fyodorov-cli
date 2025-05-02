package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api-client"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	resourceTypes = common.Enum{"models", "agents", "tools", "providers", "instances"}
	cacheMutex    sync.Mutex
	TTL           = 5 * time.Minute
	initialRun    = true
	cache         = Cache{
		TTL:                 TTL,
		TimeSinceLastUpdate: time.Now().Add(-TTL).Add(-1 * time.Minute),
		Resources:           common.CreateFyodorovConfig(),
	}
)

type Cache struct {
	TTL                 time.Duration
	TimeSinceLastUpdate time.Time
	Resources           *common.FyodorovConfig
}

func (c *Cache) IsExpired() bool {
	return time.Since(c.TimeSinceLastUpdate) > c.TTL
}

func (c *Cache) Update(resourceType *string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	fmt.Printf("%s Cache accessed\n", time.Now())
	if initialRun || NoCache { // Populate cache on first run
		fmt.Printf("Populating cache - initial run (%t) or no-cache (%t)\n", initialRun, NoCache)
		c.Resources = getResources(nil)
		c.TimeSinceLastUpdate = time.Now()
		initialRun = false
		return
	}
	if c.IsExpired() {
		fmt.Printf("%s cache expired (%s), updating\n", time.Now(), c.TimeSinceLastUpdate)
		resources := getResources(resourceType)
		c.TimeSinceLastUpdate = time.Now()
		if resourceType == nil {
			c.Resources = resources
			return
		}
		switch *resourceType {
		case "models":
			c.Resources.Models = resources.Models
		case "agents":
			c.Resources.Agents = resources.Agents
		case "tools":
			c.Resources.Tools = resources.Tools
		case "providers":
			c.Resources.Providers = resources.Providers
		case "instances":
			c.Resources.Instances = resources.Instances
		default:
			c.Resources = resources
		}
		fmt.Printf("%s Resources updated: %+v \n", c.TimeSinceLastUpdate, c.Resources)
	}
}

func init() {
	rootCmd.AddCommand(listResourcesCmd)
	rootCmd.AddCommand(removeResourcesCmd)
}

var listResourcesCmd = &cobra.Command{
	Use:     fmt.Sprintf("list (ls) [resource type: %s]", strings.Join(resourceTypes, "|")),
	Short:   "List deployed resources for a user",
	Aliases: []string{"ls"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		autocompleteResourceTypes := slices.DeleteFunc(resourceTypes, func(s string) bool {
			return slices.Contains(args, s)
		})
		return autocompleteResourceTypes, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		resources := GetAllResources()
		var bytes []byte
		var err error
		if len(args) > 0 {
			for _, resourceType := range args {
				if !resourceTypes.Contains(resourceType) {
					fmt.Printf("\033[33mInvalid resource type: %s. Valid types are: %v\033[0m\n", resourceType, resourceTypes)
					os.Exit(1)
				}
				switch resourceType {
				case "models":
					bytes, err = yaml.Marshal(resources.Models)
				case "agents":
					bytes, err = yaml.Marshal(resources.Agents)
				case "tools":
					bytes, err = yaml.Marshal(resources.Tools)
				case "providers":
					bytes, err = yaml.Marshal(resources.Providers)
				case "instances":
					bytes, err = yaml.Marshal(resources.Instances)
				default:
					fmt.Printf("\033[33mInvalid resource type: %s. Valid types are: %v\033[0m\n", resourceType, resourceTypes)
					os.Exit(1)
				}
				if err != nil {
					fmt.Printf("\033[33mError marshaling fyodorov config to yaml: %v\n\033[0m", err)
					return
				}
				// Print the YAML to stdout
				fmt.Printf("\033[36m---%s resources---:\033[0m\n", strings.Title(resourceType))
				fmt.Printf("%s\n", string(bytes))
			}
		} else {
			bytes, err = yaml.Marshal(resources)
			if err != nil {
				fmt.Printf("\033[33mError marshaling fyodorov config to yaml: %v\n\033[0m", err)
				return
			}
			// Print the YAML to stdout
			fmt.Printf("\033[36m---Resources---:\033[0m\n")
			fmt.Printf("%s\n", string(bytes))
		}
	},
}

var removeResourcesCmd = &cobra.Command{
	Use:     "remove (rm)",
	Short:   "Remove deployed resources for a user",
	Aliases: []string{"rm"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return resourceTypes, cobra.ShellCompDirectiveNoFileComp
		} else if len(args) >= 1 {
			resourceType := args[0]
			if !resourceTypes.Contains(resourceType) {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			// Get the list of resources for the specified type
			resources := GetResources(&resourceType)
			var resourceNames []string
			switch resourceType {
			case "models":
				for _, resource := range resources.Models {
					resourceNames = append(resourceNames, resource.String())
				}
			case "agents":
				for _, resource := range resources.Agents {
					resourceNames = append(resourceNames, resource.String())
				}
			case "tools":
				for _, resource := range resources.Tools {
					resourceNames = append(resourceNames, resource.String())
				}
			case "providers":
				for _, resource := range resources.Providers {
					resourceNames = append(resourceNames, resource.String())
				}
			case "instances":
				for _, resource := range resources.Instances {
					resourceNames = append(resourceNames, resource.String())
				}
			default:
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			for _, resource := range args[1:] {
				resourceNames = slices.DeleteFunc(resourceNames, func(s string) bool {
					return strings.Contains(s, resource)
				})
			}
			resourceNames = slices.DeleteFunc(resourceNames, func(s string) bool {
				return slices.Contains(args[1:], s)
			})
			return resourceNames, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		var resourceType string
		if len(args) < 2 {
			return
		} else if len(args) >= 2 {
			resourceType = args[0]
			if !resourceTypes.Contains(resourceType) {
				fmt.Printf("\033[33mInvalid resource type: %s. Valid types are: %v\033[0m\n", resourceType, resourceTypes)
				os.Exit(1)
			}
		}

		resourceHandles := args[1:]
		resources := GetAllResources()
		for _, resourceHandle := range resourceHandles {
			resourceId := GetResourceIDByString(resources, resourceType, resourceHandle)
			if resourceId < 1 {
				fmt.Printf("\033[33mUnable to find resource ID %s.\033[0m\n", resourceHandle)
				os.Exit(1)
			}
			DeleteResource(resourceType, resourceId)
		}
	},
}

func GetResourceIDByString(resources *common.FyodorovConfig, resourceType string, resourceString string) int64 {
	switch resourceType {
	case "models":
		for _, resource := range resources.Models {
			if resource.String() == resourceString {
				return resource.ID
			}
		}
	case "agents":
		for _, resource := range resources.Agents {
			if resource.String() == resourceString {
				return resource.ID
			}
		}
	case "tools":
		for _, resource := range resources.Tools {
			if resource.String() == resourceString {
				return resource.ID
			}
		}
	case "providers":
		for _, resource := range resources.Providers {
			if resource.String() == resourceString {
				return resource.ID
			}
		}
	case "instances":
		for _, resource := range resources.Instances {
			if resource.String() == resourceString {
				return resource.ID
			}
		}
	default:
		fmt.Printf("\033[33mInvalid resource type: %s. Valid types are: %v\033[0m\n", resourceType, resourceTypes)
		os.Exit(1)
	}
	return -1

}

func DeleteResources(resourceType string, resources []common.BaseModel) error {
	for _, resource := range resources {
		DeleteResource(resourceType, resource.GetID())
	}
	return nil
}

func DeleteResource(resourceType string, resourceId int64) {
	var client *api.APIClient
	config := &common.Config{
		Email:    v.GetString("email"),
		Password: v.GetString("password"),
	}
	if resourceType == "tools" {
		if !v.IsSet("tsiolkovsky-url") {
			fmt.Println("\033[33mTsiolkovsky URL is not set in config\033[0m")
			return
		}
		client = api.NewAPIClient(config, v.GetString("tsiolkovsky-url"))
	} else {
		client = api.NewAPIClient(config, v.GetString("gagarin-url"))
	}
	err := client.Authenticate()
	if err != nil {
		fmt.Println(err)
		fmt.Println("\033[33mUnable to authenticate with this config\033[0m")
		return
	}
	res, err := client.CallAPI("DELETE", fmt.Sprintf("/%s/%d", resourceType, resourceId), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[33mError deleting resource %s/%d: %v\n\033[0m", resourceType, resourceId, err)
		return
	}
	defer res.Close()
	body, err := io.ReadAll(res)
	if err != nil {
		fmt.Printf("Error reading response body while deleting resource %s/%d: %v\n", resourceType, resourceId, err)
		return
	}
	var bodyResponse bool
	err = json.Unmarshal(body, &bodyResponse)
	if err != nil {
		fmt.Printf(
			"\033[33mError unmarshaling response body after deleting resource %s/%d: %s\n\t%v\n\033[0m",
			resourceType,
			resourceId,
			string(body),
			err,
		)
		return
	}
	if !bodyResponse {
		fmt.Printf("\033[33mFailed to delete resource %s/%d\n\033[0m", resourceType, resourceId)
		return
	}
	fmt.Printf("\033[32mResource %s/%d deleted successfully\033[0m\n", resourceType, resourceId)
}

func GetAllResources() *common.FyodorovConfig {
	return GetResources(nil)
}

func GetResources(resourceType *string) *common.FyodorovConfig {
	cache.Update(resourceType)
	return cache.Resources
}

func getResources(resourceType *string) *common.FyodorovConfig {
	config := &common.Config{
		Email:    v.GetString("email"),
		Password: v.GetString("password"),
	}
	var client *api.APIClient
	if resourceType != nil && *resourceType == "tools" {
		client = api.NewAPIClient(config, v.GetString("tsiolkovsky-url"))
	} else {
		client = api.NewAPIClient(config, v.GetString("gagarin-url"))
	}
	err := client.Authenticate()
	if err != nil {
		fmt.Println(err)
		fmt.Println("\033[33mUnable to authenticate with this config\033[0m")
		os.Exit(1)
	}
	fyodorovConfig, err := client.GetResources(resourceType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[33mError fetching resources: %v\n\033[0m", err)
		os.Exit(1)
	}
	return &fyodorovConfig
}
