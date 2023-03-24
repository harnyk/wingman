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
	OS    string
	Shell string
	User  string
}

func NewContext() (EnvironmentContext, error) {
	user, err := getUser()
	if err != nil {
		return EnvironmentContext{}, err
	}

	envContext := EnvironmentContext{
		OS:    getOS(),
		Shell: getShell(),
		User:  user,
	}

	return envContext, nil
}

func getOS() string {
	return runtime.GOOS
}

func getShell() string {
	return os.Getenv("SHELL")
}

func getUser() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.Username, nil
}
