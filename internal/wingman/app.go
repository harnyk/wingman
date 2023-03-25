package wingman

import (
	"os"

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
}

func (a *App) Loop(query string) error {

	stopSpinner := StartSpinner()
	resp, err := GetResponse(a.OpenAIClient, query)
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
		return RunCommand(resp.Command)
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
