package command

import "github.com/mattherman/mhgit/objects"

// HashObject will hash an existing file and write it to the object store
// if desired.
func HashObject(filename string, write bool) string {
	return objects.HashFile(filename, write)
}
