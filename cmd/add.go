package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [files...]",
	Short: "Add file contents to the index",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addFiles(args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addFiles(files []string) {
	for _, file := range files {
		updateIndex(file, true, false)
	}
}
