package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type object struct {
	data       []byte
	objectType string
}

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
		hash := hashObject(object{data: []byte("what is up, doc?"), objectType: "blob"}, true)
		fmt.Println(hash)
	case "read":
		obj, err := readObject(args[1])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n", obj.data)
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

func hashObject(objectToHash object, write bool) string {
	header := []byte(objectToHash.objectType + " " + strconv.Itoa(len(objectToHash.data)) + "\000")

	fullData := append(header, objectToHash.data...)

	sha1 := computeSha1(fullData)

	objectPath := filepath.Join("./.git/objects", sha1[:2])

	os.Mkdir(objectPath, 0700)
	fileName := filepath.Join(objectPath, sha1[2:])

	if write {
		err := writeCompressedFile(fileName, fullData)

		if err != nil {
			panic(err)
		}
	}

	return sha1
}

func computeSha1(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	hashBytes := hasher.Sum(nil)
	return fmt.Sprintf("%x\n", hashBytes)
}

func readObject(hash string) (object, error) {
	if len(hash) < 3 {
		return object{}, errors.New("Prefix provided must be at least three characters")
	}

	hashDirectory := "./.git/objects/" + hash[:2]
	hashRemainder := hash[2:]

	files, _ := filepath.Glob(hashDirectory + "/" + hashRemainder + "*")

	if len(files) == 0 {
		return object{}, errors.New("Object " + hash + " not found.")
	} else if len(files) > 1 {
		return object{}, errors.New("Found multiple matches for " + hash + ".")
	} else {
		content, err := readCompressedFile(files[0])

		if err != nil {
			panic(err)
		}

		return object{data: content, objectType: "blob"}, nil
	}
}

func writeCompressedFile(filename string, uncompressedData []byte) error {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(uncompressedData)
	w.Close()

	return ioutil.WriteFile(filename, buf.Bytes(), 0700)
}

func readCompressedFile(filename string) ([]byte, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewReader(fileContent)
	r, err := zlib.NewReader(buf)
	defer r.Close()
	if err != nil {
		panic(err)
	}

	result, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return result, nil
}
