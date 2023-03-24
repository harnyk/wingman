package wingman

import "strings"

func CreatePrompt(userPrompt string, envContext EnvironmentContext) string {
	sbPrompt := strings.Builder{}

	sbPrompt.WriteString("Provide a terminal command which would do the following:\n")
	sbPrompt.WriteString(userPrompt)
	sbPrompt.WriteString("\n\n")
	sbPrompt.WriteString("Environment Context:\n")

	sbPrompt.WriteString("OS: ")
	sbPrompt.WriteString(envContext.OS)
	sbPrompt.WriteString("\n")

	sbPrompt.WriteString("Shell: ")
	sbPrompt.WriteString(envContext.Shell)
	sbPrompt.WriteString("\n")

	sbPrompt.WriteString("User: ")
	sbPrompt.WriteString(envContext.User)
	sbPrompt.WriteString("\n")

	sbPrompt.WriteString("Instruction: it is very important to reply in the following format (the response must terminate with END):\n\n")
	sbPrompt.WriteString("COMMAND:\n")
	sbPrompt.WriteString("some command, e.g. ls -la\n")
	sbPrompt.WriteString("DESCRIPTION:\n")
	sbPrompt.WriteString("a brief explanation how the command works\n")
	sbPrompt.WriteString("END\n")

	return sbPrompt.String()
}
