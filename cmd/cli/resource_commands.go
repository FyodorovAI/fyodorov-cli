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
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"gopkg.in/yaml.v3"
)

var (
	resourceTypes = common.Enum{"models", "agents", "tools", "providers", "instances"}
	cache         = &Cache{
		Resources: common.CreateFyodorovConfig(v),
		FileName:  common.GetPlatformBasePath() + "/cache.yaml",
		Mutex:     sync.Mutex{},
	}
)

func init() {
	cache.ReadCacheFromFile()
	rootCmd.AddCommand(listResourcesCmd)
	rootCmd.AddCommand(removeResourcesCmd)
}

type Cache struct {
	Viper     *viper.Viper
	Resources *common.FyodorovConfig
	FileName  string
	Mutex     sync.Mutex
}

func (cache *Cache) Update(forceUpdate bool) {
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()
	cache.ReadCacheFromFile()
	if cache.Resources.IsExpired(v) || NoCache || forceUpdate {
		cache.Resources = getResources(nil)
		t := time.Now()
		cache.Resources.TimeOfLastCacheUpdate = &t
		yamlBytes, err := yaml.Marshal(cache.Resources)
		if err != nil {
			fmt.Printf("Error marshaling fyodorov config to yaml for cache: %v\n", err)
			return
		}
		err = os.WriteFile(cache.FileName, yamlBytes, 0644)
		if err != nil {
			fmt.Printf("Error writing cache file (%s): %v\n", cache.FileName, err)
			return
		}
	}
}

func (cache *Cache) ReadCacheFromFile() error {
	fileBytes, err := os.ReadFile(cache.FileName)
	if err != nil {
		fmt.Printf("Error reading cache file (%s): %v\n", cache.FileName, err)
		return err
	}
	err = yaml.Unmarshal(fileBytes, &cache.Resources)
	if err != nil {
		fmt.Printf("Error unmarshaling cache: %v\n", err)
		return err
	}
	return nil
}

var listResourcesCmd = &cobra.Command{
	Use:     fmt.Sprintf("list (ls) [resource type:%s]", strings.Join(resourceTypes, "|")),
	Short:   "List deployed resources for a user",
	Aliases: []string{"ls"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		autocompleteResourceTypes := slices.DeleteFunc(resourceTypes, func(s string) bool {
			return slices.Contains(args, s)
		})
		return autocompleteResourceTypes, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		resources := GetResources()
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
			resources.TimeOfLastCacheUpdate = nil
			resources.Version = nil
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
			resources := GetResources()
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
		resources := GetResources()
		resourceIDs := []common.BaseModel{}
		for _, resourceHandle := range resourceHandles {
			resourceId := GetResourceIDByString(resources, resourceType, resourceHandle)
			if resourceId < 1 {
				fmt.Printf("\033[33mUnable to find resource ID %s.\033[0m\n", resourceHandle)
				continue
			}
			resourceIDs = append(
				resourceIDs,
				&common.Resource{
					ID: resourceId,
				},
			)
		}
		DeleteResources(resourceType, resourceIDs)
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
	cache.Update(true)
	return nil
}

func DeleteResource(resourceType string, resourceId int64) {
	var client *api.APIClient
	var err error
	if resourceType == "tools" {
		if !v.IsSet("tsiolkovsky-url") {
			fmt.Println("\033[33mTsiolkovsky URL is not set in config\033[0m")
			return
		}
		client, err = api.NewAPIClient(v, v.GetString("tsiolkovsky-url"))
	} else {
		client, err = api.NewAPIClient(v, v.GetString("gagarin-url"))
	}
	err = client.Authenticate()
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

func GetResources() *common.FyodorovConfig {
	cache.Update(false)
	return cache.Resources
}

func getResources(resourceType *string) *common.FyodorovConfig {
	client, err := api.NewAPIClient(v, "")
	if err != nil {
		fmt.Println(err)
		fmt.Println("\033[33mUnable to create API client\033[0m")
		os.Exit(1)
	}
	err = client.Authenticate()
	if err != nil {
		fmt.Println(err)
		fmt.Println("\033[33mUnable to authenticate with this config\033[0m")
		os.Exit(1)
	}
	fyodorovConfig, err := client.GetResources(resourceType, v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[33mError fetching resources: %v\n\033[0m", err)
		os.Exit(1)
	}
	return fyodorovConfig
}
