package wingman

import (
	"strings"
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
