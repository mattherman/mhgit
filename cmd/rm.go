package cmd

import (
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [files...]",
	Short: "Remove files from the working tree and from the index",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeFiles(args)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

func removeFiles(files []string) {
	for _, file := range files {
		updateIndex(file, false, true)
	}
}
