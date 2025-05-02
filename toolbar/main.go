package main

import (
	"fmt"
	"os"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/api-client"
	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/spf13/viper"
)

var (
	localModelsEnabled bool
	v                  *viper.Viper
)

func main() {
	fmt.Println("Starting toolbar...")
	v = common.InitViper()
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("Fyodorov AI Toolbar")
	mEnableLocalModels := systray.AddMenuItemCheckbox("Enable local models", "Enable local models", localModelsEnabled)
	systray.AddSeparator()
	mSettings := systray.AddMenuItem("Settings", "Open the settings window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(icon.Data)

	// Handle menu item clicks
	go func() {
		for {
			select {
			case <-mSettings.ClickedCh:
				// Open the settings window here
				openSettings()
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mEnableLocalModels.ClickedCh:
				// Update the menu item
				if mEnableLocalModels.Checked() {
					mEnableLocalModels.Uncheck()
					localModelsEnabled = false
				} else {
					mEnableLocalModels.Check()
					localModelsEnabled = true
				}
				enableLocalModels()
			}
		}
	}()

	authenticate()
}

func authenticate() {
	config, err := common.GetConfig(nil, v)
	if err != nil {
		fmt.Println("No config file found")
	}
	fmt.Println("Found config file")

	if config != nil && config.Email != "" && config.Password != "" {
		client := api.NewAPIClient(config, config.GagarinURL)
		err = client.Authenticate()
		if err != nil {
			fmt.Println("Error authenticating:", err)
			return
		} else {
			fmt.Println("Authenticated successfully")
		}
	} else {
		// download latest executable from github releases based platform
		repo := "FyodorovAI/fyodorov-cli"
		assets, err := fetchLatestRelease(repo)
		if err != nil {
			fmt.Printf("Error fetching latest release: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Fetched latest release assets", assets) // temp log line
		path, err := downloadAppropriateAsset(assets)
		if err != nil {
			fmt.Printf("Error downloading asset: %v\n", err)
			os.Exit(1)
		}

		openTerminalWithCommand(fmt.Sprintf("%s auth", path))
	}
}

func openSettings() {
	fmt.Println("Open settings")
}

func enableLocalModels() {
	fmt.Println("local models checked:", localModelsEnabled)
	ollama()
}

func onExit() {
	// clean up here
	fmt.Println("Exiting...")
	os.Exit(0)
}
