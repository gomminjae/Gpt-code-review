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

			// Display progress bar while waiting for GPT response
			done := make(chan bool)
			go showProgressBar(done)

			review, err := internal.ReviewCode(filePath)
			done <- true // Signal the progress bar to stop

			if err != nil {
				fmt.Printf("❌ Error: %s\n", err)
				return
			}

			fmt.Println("✅ Review Result:")
			fmt.Println(review)

			// Prompt user for git commit and push
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Do you want to commit and push the changes to Git? (y/n): ")
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(strings.ToLower(choice))

			if choice == "y" {
				// Perform git commit and push
				if err := gitCommitAndPush(filePath); err != nil {
					fmt.Printf("❌ Error during git push: %s\n", err)
				} else {
					fmt.Println("✅ Changes successfully committed and pushed to Git!")
				}
			} else {
				fmt.Println("❌ Skipping git push.")
			}
		},
	}
}

// showProgressBar displays a simple progress bar
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

// gitCommitAndPush handles the git commit and push commands
func gitCommitAndPush(filePath string) error {
	// Stage the file
	stageCmd := exec.Command("git", "add", filePath)
	stageCmd.Stdout = os.Stdout
	stageCmd.Stderr = os.Stderr
	if err := stageCmd.Run(); err != nil {
		return fmt.Errorf("failed to stage file: %w", err)
	}

	// Commit the changes
	commitMessage := "Code reviewed by GPT-CLI"
	commitCmd := exec.Command("git", "commit", "-m", commitMessage)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	// Push the changes
	pushCmd := exec.Command("git", "push")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}

	return nil
}
