package command

import (
	"fmt"
	"os"
	"path"

	"github.com/mattherman/mhgit/utils"
)

// InitializeRepo will create an empty repository in the current directory
// or return an error if one already exists.
func InitializeRepo(directory string) {
	gitDir := path.Join(directory, ".git")
	if !utils.FileDoesNotExist(gitDir) {
		fmt.Println("A git repository already exists in this directory")
		return
	}

	err := createInitialDirectoriesAndFiles(gitDir)
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
	f.Close()

	return err
}
