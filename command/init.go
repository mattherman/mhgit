package command

import (
	"fmt"
	"os"
)

// InitializeRepo will create an empty repository in the current directory
// or return an error if one already exists.
func InitializeRepo() {
	if !fileDoesNotExist("./.git") {
		fmt.Println("A git repository already exists in this directory")
		return
	}

	err := createInitialDirectoriesAndFiles()
	if err != nil {
		fmt.Println("Unable to create necessary files in .git directory.")
	} else {
		fmt.Println("Initialized empty Git repository.")
	}
}

func createInitialDirectoriesAndFiles() error {
	os.Mkdir("./.git", 0700)
	os.Mkdir("./.git/objects", 0700)
	os.Mkdir("./.git/refs", 0700)
	os.Mkdir("./.git/refs/heads", 0700)

	f, err := os.Create("./.git/HEAD")
	f.Close()

	return err
}

func fileDoesNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
