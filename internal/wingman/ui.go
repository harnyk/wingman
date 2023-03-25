package wingman

import (
	"fmt"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/alecthomas/chroma/quick"
	"github.com/fatih/color"
	"github.com/janeczku/go-spinner"
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

func DisplayResponse(query string, response Response) {

	formattedExplanation := string(markdown.Render(response.Explanation, 80, 0))

	var formattedCommand string
	sb := &strings.Builder{}
	err := quick.Highlight(sb, response.Command, "bash", "terminal16m", "monokai")
	if err != nil {
		formattedCommand = response.Command
	} else {
		formattedCommand = sb.String()
	}

	headingColor := color.New(color.FgHiYellow, color.Bold, color.Underline)

	headingColor.Println("             Query              ")
	fmt.Println()
	fmt.Println(query)
	fmt.Println()
	headingColor.Println("            Command             ")
	fmt.Println()
	fmt.Println(formattedCommand)
	fmt.Println()
	headingColor.Println("          Explanation           ")
	fmt.Println()
	fmt.Println(formattedExplanation)
	fmt.Println()
}

func ReviseQuery(initialQuery string) (string, error) {
	prompt := promptui.Prompt{
		Label:     "Query",
		Default:   initialQuery,
		AllowEdit: true,
	}
	query, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed %v", err)
	}

	return query, nil
}

func StartSpinner() func() {
	sp := spinner.StartNew("Loading...")
	sp.SetCharset([]string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"})

	return func() {
		sp.Stop()
	}
}
