package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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
		hashObject("abcdefg", "blob")
	}
}

func initializeRepo() {
	if fileDoesNotExist("./.git") {
		createInitialDirectoriesAndFiles()
		fmt.Println("Initialized empty repository.")
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

func hashObject(data string, objectType string) {
	header := []byte(objectType + " " + strconv.Itoa(len(data)))
	fmt.Println(header)

	fullData := append(header, 0)
	fullData = append(fullData, data...)
	fmt.Println(fullData)

	sha1 := computeSha1(fullData)
	fmt.Println(sha1)

	objectPath := "./.git/objects/" + sha1[:2]
	fmt.Println(objectPath)

	os.Mkdir(objectPath, 0700)
	fileName := objectPath + "/" + sha1[2:]
	fmt.Println(fileName)

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	defer w.Close()
	w.Write(fullData)

	bufBytes := buf.Bytes()
	fmt.Println(bufBytes)
	ioutil.WriteFile(fileName, buf.Bytes(), 0700)
}

func computeSha1(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
