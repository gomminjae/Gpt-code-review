package cmd

import (
	"fmt"
	"gpt-code-review/internal"
	"github.com/spf13/cobra"
)

func NewCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "review", // 기존 "check"를 "review"로 변경
		Short: "Run a GPT-based code review on a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filePath := args[0]
			fmt.Printf("Reviewing file: %s\n", filePath)

			review, err := internal.ReviewCode(filePath)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}

			fmt.Println("Review Result:")
			fmt.Println(review)
		},
	}
}
