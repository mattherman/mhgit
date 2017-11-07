package command

import (
	"fmt"

	"github.com/mattherman/mhgit/objects"
)

// HashObject will hash an existing file and write it to the object store
// if desired.
func HashObject(filename string, write bool) {
	hash := objects.HashFile(filename, write)
	fmt.Println(hash)
}
