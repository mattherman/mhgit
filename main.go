package main

import (
	"fmt"
	"os"

	"github.com/mattherman/mhgit/objects"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Too few arguments provided.")
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		initializeRepo()
	case "hash":
		hash := objects.HashObject(objects.Object{Data: []byte("what is up, doc?"), ObjectType: "blob"}, true)
		fmt.Println(hash)
	case "read":
		obj, err := objects.ReadObject(args[1])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s : %s\n", obj.ObjectType, obj.Data)
	}
}

func initializeRepo() {
	if fileDoesNotExist("./.git") {
		createInitialDirectoriesAndFiles()
		fmt.Println("Initialized empty Git repository.")
	} else {
		fmt.Println("A git repository already exists in this directory.")
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
