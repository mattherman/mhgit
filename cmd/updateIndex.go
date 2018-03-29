package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// updateIndexCmd represents the updateIndex command
var updateIndexCmd = &cobra.Command{
	Use:   "update-index [file]",
	Short: "Register file contents in the working tree to the index.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("updateIndex called")
	},
}

func init() {
	rootCmd.AddCommand(updateIndexCmd)
}
