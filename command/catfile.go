package command

import (
	"fmt"

	"github.com/mattherman/mhgit/objects"
)

// CatFile will inspect a stored Git object or return an error if it
// cannot be found.
func CatFile(objectName string, outputObject bool, outputType bool, outputSize bool) {
	obj, err := objects.ReadObject(objectName)

	if err != nil {
		fmt.Println(err)
	} else {
		if outputType {
			fmt.Println(obj.ObjectType)
		} else if outputSize {
			fmt.Printf("%d\n", len(obj.Data))
		} else {
			fmt.Printf("%s", obj.Data)
		}
	}
}
