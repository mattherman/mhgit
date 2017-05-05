package command

import (
	"errors"
	"os"

	"github.com/mattherman/mhgit/objects"
)

// InitializeRepo will create an empty repository in the current directory
// or return an error if one already exists.
func InitializeRepo() error {
	if !fileDoesNotExist("./.git") {
		return errors.New("A git repository already exists in this directory")
	}

	createInitialDirectoriesAndFiles()
	return nil
}

// HashObject will hash an existing file and write it to the object store
// if desired.
func HashObject(filename string, write bool) string {
	return objects.HashFile(filename, write)
}

// CatFile will inspect a stored Git object or return an error if it
// cannot be found.
func CatFile(objectName string) (objects.Object, error) {
	return objects.ReadObject(objectName)
}

func createInitialDirectoriesAndFiles() {
	os.Mkdir("./.git", 0700)
	os.Mkdir("./.git/objects", 0700)
	os.Mkdir("./.git/refs", 0700)
	os.Mkdir("./.git/refs/heads", 0700)

	f, _ := os.Create("./.git/HEAD")
	f.Close()
}

func fileDoesNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
