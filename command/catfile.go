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
			fmt.Printf("%v", obj.ObjectType)
		} else if outputSize {
			fmt.Printf("%d", len(obj.Data))
		} else {
			prettyPrint(obj)
		}
	}
}

func prettyPrint(obj objects.Object) {
	switch obj.ObjectType {
	case "blob":
		fmt.Printf("%s", obj.Data)
	case "tree":
		fmt.Print("Its a tree!")
	case "commit":
		fmt.Print("Commmmmmiiiiiittt")
	}
}
