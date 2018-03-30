package cmd

import (
	"fmt"

	"github.com/mattherman/mhgit/index"
	"github.com/spf13/cobra"
)

// lsFilesCmd represents the lsFiles command
var lsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "Show information about files in the index and the working tree",
	Run: func(cmd *cobra.Command, args []string) {
		index, err := index.ReadIndex()
		if err != nil {
			fmt.Printf("Could not read index: %v\n", err)
		}

		for _, entry := range index.Entries {
			if showStaged {
				fmt.Printf("%o %s 0\t %s\n", entry.Mode, entry.Hash, entry.Path)
			} else {
				fmt.Printf("%s\n", entry.Path)
			}
		}
	},
}

var showStaged bool

func init() {
	rootCmd.AddCommand(lsFilesCmd)
	lsFilesCmd.Flags().BoolVarP(&showStaged, "stage", "s", false, "Show staged contents' mode bits, object name, and stage number in the output.")
}
