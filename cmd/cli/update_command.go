package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

var (
	version = "0.4.2"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

// GitHubRelease represents the GitHub release structure for the latest release.
type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

// Fyodorov commands
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Fyodorov command line tool",
	Run: func(cmd *cobra.Command, args []string) {
		// Replace 'owner' and 'repo' with the actual owner and repository
		url := "https://api.github.com/repos/FyodorovAI/fyodorov-cli/releases/latest"

		// Send GET request to GitHub API
		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			return
		}
		defer response.Body.Close()

		// Read the body of the response
		data, _ := io.ReadAll(response.Body)

		if response.StatusCode != 200 {
			fmt.Println("Unable to fetch the latest release")
			fmt.Printf("Response status: %s\n", response.Status)
			fmt.Printf("Response body: %s\n", string(data))
			return
		}

		var release GitHubRelease
		if err := json.Unmarshal(data, &release); err != nil {
			fmt.Printf("Error parsing the response body: %s\n", err)
			return
		}

		// Print the latest release version number
		fmt.Printf("Latest release version: %s\n", release.TagName)

		localVersion, err := semver.NewVersion(version)
		if err != nil {
			fmt.Printf("Error parsing local version %s: %s\n", version, err)
			return
		}

		remoteVersion, err := semver.NewVersion(release.TagName)
		if err != nil {
			fmt.Printf("Error parsing github remote version %s: %s\n", release.TagName, err)
			return
		}

		// Compare the versions
		if localVersion.LessThan(remoteVersion) {
			fmt.Printf(
				"There is an update available(%s): %s\n",
				release.TagName,
				"https://github.com/FyodorovAI/fyodorov-cli/releases/tag/"+release.TagName,
			)
		} else {
			fmt.Printf("You have the latest version: %s\n", version)
		}
	},
}
