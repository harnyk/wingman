package wingman

import (
	"fmt"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/manifoldco/promptui"
)

// use promptui to create a menu:
// 1. Run This Command
// 2. Revise Query
// 3. Exit

type MenuAction int

const (
	MARunCommand MenuAction = 1 + iota
	MAReviseQuery
	MAExit
)

func Menu() (MenuAction, error) {
	prompt := promptui.Select{
		Label: "Select Action",
		Items: []string{"Run This Command", "Revise Query", "Exit"},
	}
	index, _, err := prompt.Run()

	if err != nil {
		return 0, fmt.Errorf("Prompt failed %v", err)
	}

	return MenuAction(index + 1), nil
}

func DisplayResponse(response Response) {

	renderedExplanation := string(markdown.Render(response.Explanation, 80, 10))

	fmt.Println("________________Command___________________")
	fmt.Println()
	fmt.Println(response.Command)
	fmt.Println("______________Explanation_________________")
	fmt.Println()
	fmt.Println(renderedExplanation)
	fmt.Println()
}

func ReviseQuery(initialQuery string) (string, error) {
	prompt := promptui.Prompt{
		Label: "Query",
	}
	query, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("Prompt failed %v", err)
	}

	return query, nil
}
