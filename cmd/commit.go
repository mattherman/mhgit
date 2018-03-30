package cmd

import (
	"fmt"

	"github.com/mattherman/mhgit/objects"
	"github.com/mattherman/mhgit/refs"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Run: func(cmd *cobra.Command, args []string) {
		commit(commitMsg)
	},
}

var commitMsg string

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&commitMsg, "message", "m", "", "The commit message")
	commitCmd.MarkFlagRequired("message")
}

func commit(message string) {
	treeHash, err := writeTree()
	if err != nil {
		fmt.Printf("Failed to write the tree object to the database: %v\n", err)
	}

	latestCommit, _ := refs.LatestCommit()
	parent := ""
	if latestCommit != "" {
		parent = fmt.Sprintf("parent %s\n", latestCommit)
	}
	author := "Matthew Herman <mattherman11@gmail.com> 1493552135 -0500"
	committer := "committer Matthew Herman <mattherman11@gmail.com> 1493552135 -0500"

	fullCommit := fmt.Sprintf("tree %s\n%sauthor %s\ncommitter %s\n\n%s\n", treeHash, parent, author, committer, message)
	fullCommitBytes := []byte(fullCommit)

	obj := objects.Object{ObjectType: "commit", Data: fullCommitBytes}
	hash, err := objects.HashObject(obj, true)
	if err != nil {
		fmt.Printf("Failed to write the commit to the database: %v\n", err)
	} else {
		err := refs.UpdateLatestCommit(hash)
		if err != nil {
			fmt.Printf("Failed to update latest commit: %v\n", err)
		}
		fmt.Printf("%s\n", hash)
	}
}
