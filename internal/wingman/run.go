package wingman

import (
	"log"
	"os"
	"os/exec"
)

func RunCommand(command string) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
		log.Printf("No SHELL environment variable found, using %s", shell)
	}
	cmd := exec.Command(shell, "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
