package cmd

import (
	"fmt"

	"github.com/mattherman/mhgit/refs"
	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			createBranch(args[0])
		} else {
			listBranches()
		}

	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}

func createBranch(branchName string) {
	err := refs.CreateBranch(branchName)
	if err != nil {
		fmt.Printf("Failed to create branch: %v\n", err)
	}
}

func listBranches() {
	currentBranch, err := refs.CurrentBranch()
	if err != nil {
		fmt.Printf("Failed to lookup the current branch: %v\n", err)
		return
	}

	branches, err := refs.ListBranches()
	if err != nil {
		fmt.Printf("Failed to retrieve branches: %v\n", err)
		return
	}

	for _, branch := range branches {
		if branch == currentBranch {
			fmt.Printf("* %s\n", branch)
		} else {
			fmt.Printf("  %s\n", branch)
		}
	}
}
