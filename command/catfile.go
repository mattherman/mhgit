package command

import (
	"fmt"

	"github.com/mattherman/mhgit/objects"
)

// CatFile will inspect a stored Git object or return an error if it
// cannot be found.
func CatFile(objectName string) {
	obj, err := objects.ReadObject(objectName)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s : %s\n", obj.ObjectType, obj.Data)
	}
}
