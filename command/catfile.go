package command

import (
	"github.com/mattherman/mhgit/objects"
)

// CatFile will inspect a stored Git object or return an error if it
// cannot be found.
func CatFile(objectName string) (objects.Object, error) {
	return objects.ReadObject(objectName)
}
