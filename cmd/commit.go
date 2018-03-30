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
		err := commit(commitMsg)
		if err != nil {
			fmt.Printf("Failed to commit the changes: %v\n", err)
		}
	},
}

var commitMsg string

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&commitMsg, "message", "m", "", "The commit message")
	commitCmd.MarkFlagRequired("message")
}

func commit(message string) error {
	treeHash, err := writeTree()
	if err != nil {
		return err
	}

	// TODO please god fix this...if can't find the file must be first commit hahaha i hate myself
	latestCommit, _ := refs.LatestCommit()
	parent := ""
	if latestCommit != "" {
		parent = fmt.Sprintf("parent %s\n", latestCommit)
	}

	// TODO do not hardcode this junk
	author := "Matthew Herman <mattherman11@gmail.com> 1493552135 -0500"
	committer := "committer Matthew Herman <mattherman11@gmail.com> 1493552135 -0500"

	fullCommit := fmt.Sprintf("tree %s\n%sauthor %s\ncommitter %s\n\n%s\n", treeHash, parent, author, committer, message)
	fullCommitBytes := []byte(fullCommit)

	obj := objects.Object{ObjectType: "commit", Data: fullCommitBytes}
	hash, err := objects.HashObject(obj, true)
	if err != nil {
		return err
	}

	err = refs.UpdateLatestCommit(hash)
	return err
}
