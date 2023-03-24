package wingman

import (
	"os"
	"os/exec"
)

func RunCommand(command string) error {
	// todo: make this work on windows too
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
