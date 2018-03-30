package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattherman/mhgit/index"
	"github.com/mattherman/mhgit/objects"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Run: func(cmd *cobra.Command, args []string) {
		showStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func showStatus() {
	workingDir, err := getWorkingDirectoryFiles()
	if err != nil {
		fmt.Printf("Failed to retrieve working directory: %v\n", err)
	}

	index, err := index.ReadIndex()
	if err != nil {
		fmt.Printf("Failed to retrieve indexed files: %v\n", err)
	}

	status := getStatus(index.Entries, workingDir)
	printStatus(status)
}

func getWorkingDirectoryFiles() ([]string, error) {
	var files []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {

		// Do not add directories to collection and enforce exclusion of '.git'
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}

type status struct {
	added    []string
	modified []string
	removed  []string
}

type hashPair struct {
	indexHash      string
	workingDirHash string
}

func getStatus(indexEntries []index.Entry, workingDir []string) status {
	statusMap := make(map[string]hashPair)
	var added []string
	var modified []string
	var removed []string

	for _, entry := range indexEntries {
		statusMap[entry.Path] = hashPair{indexHash: entry.Hash}
	}
	for _, path := range workingDir {
		hash, _ := objects.HashFile(path, false)
		existing := statusMap[path]
		existing.workingDirHash = hash
		statusMap[path] = existing
	}
	for k, v := range statusMap {
		if v.indexHash != "" && v.workingDirHash != "" && v.indexHash != v.workingDirHash {
			modified = append(modified, k)
		} else if v.indexHash == "" && v.workingDirHash != "" {
			added = append(added, k)
		} else if v.indexHash != "" && v.workingDirHash == "" {
			removed = append(removed, k)
		}
	}

	return status{added: added, modified: modified, removed: removed}
}

func printStatus(status status) {
	if len(status.modified) > 0 || len(status.removed) > 0 {
		fmt.Print("\nChanges not staged for commit:\n\n")
		for _, path := range status.modified {
			fmt.Printf("\tmodified: %s\n", path)
		}
		for _, path := range status.removed {
			fmt.Printf("\tdeleted: %s\n", path)
		}
	}

	if len(status.added) > 0 {
		fmt.Print("\nUntracked files:\n\n")
		for _, path := range status.added {
			fmt.Printf("\t%s\n", path)
		}
	}

	fmt.Println()
}
