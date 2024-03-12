package main

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

type UnixNotifier struct{}

func (n UnixNotifier) Notify(title, message string) error {
	// if runtime.GOOS == "windows" {
	// 	notification := toast.Notification{
	// 		AppID:   "Your App",
	// 		Title:   title,
	// 		Message: message,
	// 		// You can specify more properties here
	// 	}
	// 	return notification.Push()
	// }
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// macOS uses osascript for notifications
		script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
		cmd = exec.Command("osascript", "-e", script)
	case "linux":
		// Linux typically uses notify-send
		cmd = exec.Command("notify-send", title, message)
	default:
		return errors.New("unsupported platform")
	}
	return cmd.Run()
}
