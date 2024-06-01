package wingman

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

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
	client     *openai.Client
	envContext EnvironmentContext
	model      string
}

func NewApp(openAIClient *openai.Client, openaiModel string) (*App, error) {
	context, err := NewContext()
	if err != nil {
		return nil, err
	}

	return &App{
		client:     openAIClient,
		envContext: context,
		model:      openaiModel,
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
	stopSpinner()
	if err != nil {
		return err
	}

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

	resp, err := a.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: a.model,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			N: 1,
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
	response, err := ParseResponseJSON(raw.Message.Content)

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

const promptTemplateTextJsonFormat = `
Provide a terminal command which would do the following:\n
{{ .UserPrompt }}

Environment Context:
OS: {{ .Context.OS }}
Shell: {{ .Context.Shell }}
User: {{ .Context.User }}
Instruction: it is very important that your response is in the JSON format, corresponding to the following JSON schema:

{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
		"command": {
			"description": "the command to run directly in the shell",
			"type": "string"
		},
		"explanation": {
			"description": "a brief explanation how the command works",
			"type": "string"
		}
	},
	"required": ["command", "explanation"],
	"additionalProperties": false
}

Example:

User prompt:

list all files in the current directory

Response:
{
	"command": "ls -la",
	"explanation": "ls lists all files in the current directory"
}

In the explanation field use Markdown formatting if needed.
It is also important to pay attention on producing secure syntax, e.g. using proper quotes in the command.
Do not use Markdown formatting in the command field.
`

var promptTemplate = template.Must(template.New("prompt").Parse(promptTemplateTextJsonFormat))

func (a *App) createPrompt(userPrompt string) string {
	sbPrompt := strings.Builder{}
	if err := promptTemplate.Execute(&sbPrompt, map[string]interface{}{
		"Context":    a.envContext,
		"UserPrompt": userPrompt,
	}); err != nil {
		log.Fatal(err)
	}
	return sbPrompt.String()
}
