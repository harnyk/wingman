package wingman

import (
	"fmt"
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

func (a *App) Loop(initialQuery string) error {
	resp, err := GetResponse(a.OpenAIClient, initialQuery)
	if err != nil {
		return err
	}

	DisplayResponse(resp)

	action, err := Menu()
	if err != nil {
		return err
	}

	fmt.Println("Action: ", action)

	switch action {
	case MARunCommand:
		return RunCommand(resp.Command)
	case MAReviseQuery:
		newQuery, err := ReviseQuery(initialQuery)
		if err != nil {
			return err
		}
		return a.Loop(newQuery)
	case MAExit:
		os.Exit(0)
	}

	return nil
}
