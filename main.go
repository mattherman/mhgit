package main

import (
	"fmt"
	"os"

	"github.com/mattherman/mhgit/command"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("mhgit", "An implementation of Git written in Go.")

	initialize = app.Command("init", "Initialize a new repository.")

	hashObject      = app.Command("hash-object", "Hash an existing file.")
	hashObjectWrite = hashObject.Flag("write", "Whether or not to write the object to the Git object store.").Short('w').Bool()
	hashObjectFile  = hashObject.Arg("file", "The path to the file being hashed.").Required().String()

	catFile       = app.Command("cat-file", "Inspect a stored Git object.")
	catFilePretty = catFile.Flag("pretty", "Pretty-print the object based on type.").Short('p').Bool()
	catFileType   = catFile.Flag("type", "Output the type of the object.").Short('t').Bool()
	catFileSize   = catFile.Flag("size", "Output the size of the object.").Short('s').Bool()
	catFileObject = catFile.Arg("object", "The name of the object to show.").Required().String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case initialize.FullCommand():
		err := command.InitializeRepo()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Initialized empty Git repository.")
		}
	case hashObject.FullCommand():
		hash := command.HashObject(*hashObjectFile, *hashObjectWrite)
		fmt.Println(hash)
	case catFile.FullCommand():
		obj, err := command.CatFile(*catFileObject)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s : %s\n", obj.ObjectType, obj.Data)
		}
	}
}
