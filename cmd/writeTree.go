package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/mattherman/mhgit/index"
	"github.com/mattherman/mhgit/objects"

	"github.com/spf13/cobra"
)

// writeTreeCmd represents the writeTree command
var writeTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Create a tree object from the current index",
	Run: func(cmd *cobra.Command, args []string) {
		writeTree()
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}

func writeTree() {
	index, err := index.ReadIndex()
	if err != nil {
		fmt.Printf("Failed to read index: %v\n", err)
	}

	var blobBytes []byte
	for _, entry := range index.Entries {
		hashAsBytes, _ := hex.DecodeString(entry.Hash)
		// TODO do not hardcode the file mode
		blob := fmt.Sprintf("100644 %s\000%s", entry.Path, hashAsBytes)
		blobBytes = append(blobBytes, []byte(blob)...)
	}

	obj := objects.Object{ObjectType: "tree", Data: blobBytes}

	hash, err := objects.HashObject(obj, true)
	if err != nil {
		fmt.Printf("Failed to write the tree to the database: %v\n", err)
	} else {
		fmt.Printf("%s\n", hash)
	}
}
