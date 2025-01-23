package cmd

import (
	"bufio"
	"fmt"
	"gpt-code-review/internal"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// NewCheckCommand creates the "review" command
func NewCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "review",
		Short: "Run a GPT-based code review on a file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filePath := args[0]
			fmt.Printf("Reviewing file: %s\n", filePath)

			// Show progress bar during GPT processing
			done := make(chan bool)
			go showProgressBar(done)

			review, err := internal.ReviewCode(filePath)
			close(done) // Ensure the channel is closed to terminate progress bar

			if err != nil {
				fmt.Printf("❌ Error: %s\n", err)
				return
			}

			// Display review results
			fmt.Println("✅ Review Result:")
			fmt.Println(review)

			// Ask for Git commit and push
			if askForGitCommit() {
				if err := gitCommitAndPush(filePath); err != nil {
					fmt.Printf("❌ Error during Git operation: %s\n", err)
				} else {
					fmt.Println("✅ Changes successfully committed and pushed to Git!")
				}
			} else {
				fmt.Println("❌ Skipping Git push.")
			}
		},
	}
}

// showProgressBar displays a progress bar during long-running tasks
func showProgressBar(done chan bool) {
	for {
		select {
		case <-done:
			fmt.Print("\rProcessing complete!          \n")
			return
		default:
			for _, r := range `|/-\` {
				fmt.Printf("\rProcessing... %c", r)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// askForGitCommit prompts the user for confirmation to commit and push changes
func askForGitCommit() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Do you want to commit and push the changes to Git? (y/n): ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("❌ Error reading input: %s\n", err)
			continue
		}

		choice = strings.TrimSpace(strings.ToLower(choice))
		if choice == "y" {
			return true
		} else if choice == "n" {
			return false
		} else {
			fmt.Println("❌ Invalid input. Please enter 'y' or 'n'.")
		}
	}
}

// gitCommitAndPush performs Git operations to stage, commit, and push changes
func gitCommitAndPush(filePath string) error {
	commands := [][]string{
		{"git", "add", filePath},
		{"git", "commit", "-m", "Code reviewed by GPT-CLI"},
		{"git", "push"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to execute '%s': %w", strings.Join(cmdArgs, " "), err)
		}
	}
	return nil
}
