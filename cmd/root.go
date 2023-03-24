/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/harnyk/wingman/internal/wingman"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                "wingman",
	Short:              "AI helper for terminal",
	Long:               `Wingman is a command line tool that helps you to run Linux commands using OpenAI's GPT-3 API.`,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		openAIToken := viper.GetString("openai_token")
		if openAIToken == "" {
			return fmt.Errorf("openai_token token is not set. Please set it in config file or environment variable")
		}

		log.Println("openai_token:", openAIToken)

		singleArg := strings.Join(args, " ")

		environmentContext, err := wingman.NewContext()
		if err != nil {
			return err
		}

		prompt := wingman.CreatePrompt(singleArg, environmentContext)

		log.Println("prompt: \n", prompt)

		client := openai.NewClient(openAIToken)

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)

		if err != nil {
			return err
		}

		log.Println("response:", resp)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wingman.yaml)")

	// // Cobra also supports local flags, which will only run
	// // when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".wingman" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".wingman")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
