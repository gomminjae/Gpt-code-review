package main

import (
	"gpt-code-review/cmd"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  Warning: Could not load .env file")
	}

	// Create the root command
	rootCmd := &cobra.Command{
		Use:   "gpt-code-review",
		Short: "A CLI tool for GPT-based code review",
		Long: `GPT Code Review is a CLI tool that uses OpenAI's GPT model 
to review source code files and provide feedback directly in your terminal.`,
	}

	// Register commands
	rootCmd.AddCommand(cmd.NewCheckCommand()) // Review command

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("❌ Error: %v", err)
		os.Exit(1)
	}
}
