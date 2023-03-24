package wingman

import (
	"context"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Response struct {
	Command     string
	Explanation string
}

// example:
// ai noise
// [COMMAND]:
// ls -la
// the command is multiline
// [EXPLANATION]:
// list all files in the current directory
// this is a multiline explanation
// [END]
// ai noise

func ParseResponse(response string) (Response, error) {
	lines := strings.Split(response, "\n")

	sbCommand := strings.Builder{}
	sbExplanation := strings.Builder{}

	var sbCurrent *strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "[COMMAND]:") {
			sbCurrent = &sbCommand
			continue
		}

		if strings.HasPrefix(line, "[EXPLANATION]:") {
			sbCurrent = &sbExplanation
			continue
		}

		if strings.HasPrefix(line, "[END]") {
			break
		}

		if sbCurrent == nil {
			continue
		}
		sbCurrent.WriteString(line)
		sbCurrent.WriteString("\n")
	}

	return Response{
		Command:     strings.Trim(sbCommand.String(), "\n"),
		Explanation: strings.Trim(sbExplanation.String(), "\n"),
	}, nil
}

func GetResponse(client *openai.Client, query string) (Response, error) {

	environmentContext, err := NewContext()
	if err != nil {
		return Response{}, err
	}

	prompt := CreatePrompt(query, environmentContext)

	resp, err := client.CreateChatCompletion(
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
