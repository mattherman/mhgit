package cmd

import (
	"fmt"

	"github.com/mattherman/mhgit/objects"
	"github.com/spf13/cobra"
)

// hashObjectCmd represents the hashObject command
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object [file]",
	Short: "Compute object ID and optionally creates a blob from a file.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		HashObject(args[0], write)
	},
}

var objectType string
var write bool

func init() {
	rootCmd.AddCommand(hashObjectCmd)
	hashObjectCmd.Flags().BoolVarP(&write, "write", "w", false, "Whether or not to write the object to the git object store.")
	hashObjectCmd.Flags().StringVarP(&objectType, "type", "t", "blob", "The type of object. Defaults to 'blob'.")
}

// HashObject will hash an existing file and write it to the object store
// if desired.
func HashObject(filename string, write bool) {
	hash, err := objects.HashFile(filename, write)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(hash)
	}
}
