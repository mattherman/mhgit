package command

import (
	"fmt"

	"github.com/mattherman/mhgit/objects"
)

// HashObject will hash an existing file and write it to the object store
// if desired.
func HashObject(filename string, write bool) {
	hash, err := objects.HashFile(filename, write)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(hash)
	}
}
