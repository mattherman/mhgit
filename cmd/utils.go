package cmd

import "os"

// FileDoesNotExist returns a bool with the
// value true if the file at the specified path
// does not exist, else returns false.
func fileDoesNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
