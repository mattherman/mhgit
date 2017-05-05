package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/mattherman/mhgit/objects"
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
		execInitializeRepo()
	case hashObject.FullCommand():
		execHashObject(*hashObjectFile, *hashObjectWrite)
	case catFile.FullCommand():
		execCatFile(*catFileObject)
	}
}

func execInitializeRepo() {
	if fileDoesNotExist("./.git") {
		createInitialDirectoriesAndFiles()
		fmt.Println("Initialized empty Git repository.")
	} else {
		fmt.Println("A git repository already exists in this directory.")
	}
}

func execHashObject(filename string, write bool) {
	hash := objects.HashFile(filename, write)
	fmt.Println(hash)
}

func execCatFile(objectName string) {
	obj, err := objects.ReadObject(objectName)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s : %s\n", obj.ObjectType, obj.Data)
	}
}

func createInitialDirectoriesAndFiles() {
	os.Mkdir("./.git", 0700)
	os.Mkdir("./.git/objects", 0700)
	os.Mkdir("./.git/refs", 0700)
	os.Mkdir("./.git/refs/heads", 0700)

	f, _ := os.Create("./.git/HEAD")
	f.Close()
}

func fileDoesNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
