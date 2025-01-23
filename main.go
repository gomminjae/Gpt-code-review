package main

import (
	"gpt-code-review/cmd"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	// Root Command 생성
	rootCmd := &cobra.Command{
		Use:   "test-cli",
		Short: "CLI for testing commands",
	}

	// CheckCommand 추가
	rootCmd.AddCommand(cmd.NewCheckCommand())

	// Command 실행
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(1)
	}
}
