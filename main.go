package main

import (
	"os"

	"github.com/mattherman/mhgit/command"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("mhgit", "An implementation of Git written in Go.")

	initialize          = app.Command("init", "Create an empty Git repository or reinitialize an existing one.")
	initializeDirectory = initialize.Arg("directory", "An optional directory to create for the repository.").String()

	hashObject      = app.Command("hash-object", "Compute object ID and optionally creates a blob from a file.")
	hashObjectWrite = hashObject.Flag("write", "Whether or not to write the object to the Git object store.").Short('w').Bool()
	hashObjectType  = hashObject.Flag("type", "The type of object. Defaults to 'blob'.").Default("blob").Short('t').String()
	hashObjectFile  = hashObject.Arg("file", "The path to the file being hashed.").Required().String()

	catFile       = app.Command("cat-file", "Provide content or type and size information for repository objects.")
	catFilePretty = catFile.Flag("pretty", "Pretty-print the object based on type.").Short('p').Bool()
	catFileType   = catFile.Flag("type", "Output the type of the object.").Short('t').Bool()
	catFileSize   = catFile.Flag("size", "Output the size of the object.").Short('s').Bool()
	catFileObject = catFile.Arg("object", "The hash of the object to show.").Required().String()
)

func main() {

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case initialize.FullCommand():
		command.InitializeRepo(*initializeDirectory)

	case hashObject.FullCommand():
		command.HashObject(*hashObjectFile, *hashObjectWrite)

	case catFile.FullCommand():
		command.CatFile(*catFileObject)
	}

}
