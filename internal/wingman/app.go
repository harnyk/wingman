package wingman

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// Loop is the main loop of the application
// 1. Prompt the user for a query (or use initialQuery)
// 2. Send the query to OpenAI
// 3. Display the response
// 4. Prompt the user for an action
// 5. If the user wants to run the command, run it and exit
// 6. If the user wants to revise the query, go back to step 1
// 7. If the user wants to exit, exit

type App struct {
	OpenAIClient *openai.Client
	envContext   EnvironmentContext
}

func NewApp(openAIClient *openai.Client) (*App, error) {
	context, err := NewContext()
	if err != nil {
		return nil, err
	}

	return &App{
		OpenAIClient: openAIClient,
		envContext:   context,
	}, nil
}

func (a *App) Loop(query string) error {

	if query == "" {
		var err error
		query, err = ReviseQuery(query)
		if err != nil {
			return err
		}
	}

	stopSpinner := StartSpinner()
	resp, err := a.getResponse(query)
	if err != nil {
		stopSpinner()
		return err
	}
	stopSpinner()

	DisplayResponse(query, resp)

	action, err := Menu()
	if err != nil {
		return err
	}

	switch action {
	case MARunCommand:
		return a.runCommand(resp.Command)
	case MAReviseQuery:
		newQuery, err := ReviseQuery(query)
		if err != nil {
			return err
		}
		return a.Loop(newQuery)
	case MAExit:
		os.Exit(0)
	}

	return nil
}

func (a *App) getResponse(query string) (Response, error) {

	prompt := a.createPrompt(query)

	resp, err := a.OpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			N:     1,
			Stop: []string{
				"[END]",
			},
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return Response{}, err
	}

	raw := resp.Choices[0]
	response, err := ParseResponse(raw.Message.Content)

	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func (a *App) runCommand(command string) error {
	shell := a.envContext.Shell

	if shell == "" {
		shell = "/bin/sh"
		log.Printf("No SHELL environment variable found, using %s", shell)
	}

	args := []string{
		"-c",
		command,
	}

	if a.envContext.IsWindowsStyleFlags {
		args = []string{
			"/c",
			command,
		}
	}

	cmd := exec.Command(shell, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) createPrompt(userPrompt string) string {
	sbPrompt := strings.Builder{}

	sbPrompt.WriteString("Provide a terminal command which would do the following:\n")
	sbPrompt.WriteString(userPrompt)
	sbPrompt.WriteString("\n\n")
	sbPrompt.WriteString("Environment Context:\n")

	sbPrompt.WriteString("OS: ")
	sbPrompt.WriteString(a.envContext.OS)
	sbPrompt.WriteString("\n")

	sbPrompt.WriteString("Shell: ")
	sbPrompt.WriteString(a.envContext.Shell)
	sbPrompt.WriteString("\n")

	sbPrompt.WriteString("User: ")
	sbPrompt.WriteString(a.envContext.User)
	sbPrompt.WriteString("\n")

	sbPrompt.WriteString("Instruction: it is very important to reply in the following format (the response must terminate with [END]):\n\n")
	sbPrompt.WriteString("[COMMAND]:\n")
	sbPrompt.WriteString("some command, e.g. ls -la\n")
	sbPrompt.WriteString("[EXPLANATION]:\n")
	sbPrompt.WriteString("a brief explanation how the command works\n")
	sbPrompt.WriteString("[END]\n")
	sbPrompt.WriteString("\n")
	sbPrompt.WriteString("In the EXPLANATION field use Markdown formatting if needed\n")

	return sbPrompt.String()
}
