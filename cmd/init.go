package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Create an empty Git repository or reinitialize an existing one.",
	Run: func(cmd *cobra.Command, args []string) {
		var directory string
		if len(args) > 0 {
			directory = args[0]
		}
		InitializeRepo(directory)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// InitializeRepo will create an empty repository in the current directory
// or return an error if one already exists.
func InitializeRepo(directory string) {
	gitDir := path.Join(directory, ".git")

	_, err := os.Stat(gitDir)
	if err == nil {
		fmt.Println("A git repository already exists in this directory")
		return
	}

	err = createInitialDirectoriesAndFiles(gitDir)
	if err != nil {
		fmt.Println("Unable to create necessary files in .git directory.")
	} else {
		fmt.Println("Initialized empty Git repository.")
	}
}

func createInitialDirectoriesAndFiles(gitDir string) error {
	os.MkdirAll(gitDir, 0700)
	os.MkdirAll(path.Join(gitDir, "objects"), 0700)
	os.MkdirAll(path.Join(gitDir, "refs", "heads"), 0700)

	f, err := os.Create(path.Join(gitDir, "HEAD"))
	defer f.Close()

	if err == nil {
		_, err = f.WriteString("ref: refs/heads/master")
	}

	return err
}
