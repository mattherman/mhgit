package cmd

import (
	"github.com/mattherman/mhgit/index"
	"github.com/spf13/cobra"
)

// lsFilesCmd represents the lsFiles command
var lsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "Show information about files in the index and the working tree",
	Run: func(cmd *cobra.Command, args []string) {
		index.ReadIndex(showStaged)
	},
}

var showStaged bool

func init() {
	rootCmd.AddCommand(lsFilesCmd)
	lsFilesCmd.Flags().BoolVarP(&showStaged, "stage", "s", false, "Show staged contents' mode bits, object name, and stage number in the output.")
}