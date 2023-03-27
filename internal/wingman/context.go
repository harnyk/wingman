package wingman

import (
	"os"
	"os/user"
	"runtime"
)

// sbPrompt.WriteString("Environment\n")
// sbPrompt.WriteString("OS: Linux\n")
// sbPrompt.WriteString("Shell: bash\n")
// sbPrompt.WriteString("Terminal: xterm-256color\n")
// sbPrompt.WriteString("User: harnyk\n")

type EnvironmentContext struct {
	OS                  string
	Shell               string
	User                string
	IsWindowsStyleFlags bool
}

func NewContext() (EnvironmentContext, error) {
	user, err := getUser()
	if err != nil {
		return EnvironmentContext{}, err
	}

	sh, isWindowsStyleFlags := getShell()

	envContext := EnvironmentContext{
		OS:                  getOS(),
		Shell:               sh,
		IsWindowsStyleFlags: isWindowsStyleFlags,
		User:                user,
	}

	return envContext, nil
}

func getOS() string {
	return runtime.GOOS
}

func getShell() (shell string, isWindowsStyleFlags bool) {
	if os.Getenv("PSModulePath") != "" {
		return "powershell", true
	}
	csp := os.Getenv("ComSpec")
	if csp != "" {
		return csp, true
	}
	return os.Getenv("SHELL"), false
}

func getUser() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.Username, nil
}
