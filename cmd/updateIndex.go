package cmd

import (
	"fmt"

	"github.com/mattherman/mhgit/index"
	"github.com/spf13/cobra"
)

// updateIndexCmd represents the updateIndex command
var updateIndexCmd = &cobra.Command{
	Use:   "update-index [file]",
	Short: "Register file contents in the working tree to the index.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		updateIndex(filepath, add, remove)
	},
}

var add bool
var remove bool

func init() {
	rootCmd.AddCommand(updateIndexCmd)
	updateIndexCmd.Flags().BoolVarP(&add, "add", "a", false, "If a specified file isn’t in the index already then it’s added. Default behaviour is to ignore new files.")
	updateIndexCmd.Flags().BoolVarP(&remove, "remove", "r", false, "If a specified file is in the index but is missing then it’s removed. Default behavior is to ignore removed file.")
}

func updateIndex(filepath string, add bool, remove bool) {
	if remove {
		err := index.Remove(filepath)
		if err != nil {
			fmt.Printf("Failed to remove the index entry: %v\n", err)
		}
	} else {
		err := index.Add(filepath)
		if err != nil {
			fmt.Printf("Failed to create the index entry: %v\n", err)
		}
	}
}
