package console

import (
	"os"
	"os/exec"
	"runtime"
)

// Clear очищает консоль вывода
func Clear() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	return cmd.Run()
}
