package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/FyodorovAI/fyodorov-cli-tool/internal/common"
)

func openTerminalWithCommand(command string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// For Windows, you can use 'start' to open a new command prompt.
		cmd = exec.Command("cmd", "/c", "start", "cmd", "/k", command)
	case "linux":
		// For Linux, xterm is commonly installed. You could use others like 'gnome-terminal', 'konsole', etc.
		cmd = exec.Command("xterm", "-e", command)
	case "darwin":
		// For macOS, use 'osascript' to run an AppleScript command that opens Terminal.
		cmd = exec.Command("osascript", "-e", "tell application \"Terminal\"", "-e", "activate", "-e", "do script \""+command+"\"", "-e", "end tell")
	default:
		return fmt.Errorf("unsupported platform")
	}

	// Run the command which opens a new terminal window and executes the command.
	return cmd.Start()
}

type ReleaseInfo struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func fetchLatestRelease(repo string) ([]ReleaseInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var release struct {
		Assets []ReleaseInfo `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return release.Assets, nil
}

func downloadFile(URL, filePath string) (string, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return filePath, err
}

func downloadAppropriateAsset(assets []ReleaseInfo) (string, error) {
	var assetName string
	switch runtime.GOOS {
	case "windows":
		assetName = "fyodorov.exe"
	case "darwin":
		assetName = "fyodorov-arm"
	case "linux":
		assetName = "fyodorov-amd64"
	}

	for _, asset := range assets {
		if asset.Name == assetName {
			fmt.Printf("Downloading %s...\n", asset.Name)
			path, err := downloadFile(asset.BrowserDownloadURL, fmt.Sprintf("%s/fyodorov", common.GetPlatformBasePath()))
			if runtime.GOOS != "windows" {
				err = os.Chmod(path, 0755)
			}
			return path, err
		}
	}

	return "", fmt.Errorf("compatible asset not found")
}
