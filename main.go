package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
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
		hashObject(object{data: []byte("what is up, doc?"), objectType: "blob"}, true)
	case "read":
		obj := readObject(args[1])
		fmt.Printf("%s", obj.data)
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
	fmt.Println(sha1)

	objectPath := filepath.Join("./.git/objects", sha1[:2])

	os.Mkdir(objectPath, 0700)
	fileName := filepath.Join(objectPath, sha1[2:])
	fmt.Println(fileName)

	if write {
		var buf bytes.Buffer
		w := zlib.NewWriter(&buf)
		w.Write(fullData)
		w.Close()

		ioutil.WriteFile(fileName, buf.Bytes(), 0700)
	}

	return sha1
}

func computeSha1(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	hashBytes := hasher.Sum(nil)
	return fmt.Sprintf("%x\n", hashBytes)
}

func readObject(hash string) object {
	if len(hash) < 3 {
		fmt.Println("Prefix provided must be at least three characters.")
		return object{}
	}

	hashDirectory := "./.git/objects/" + hash[:2]
	hashRemainder := hash[2:]

	if fileDoesNotExist(hashDirectory) {
		fmt.Println("Object " + hash + " not found.")
		return object{}
	}

	files, _ := filepath.Glob(hashDirectory + "/" + hashRemainder + "*")
	if len(files) == 0 {
		fmt.Println("Object " + hash + " not found.")
	} else if len(files) > 1 {
		fmt.Println("Found multiple matches for " + hash + ".")
	} else {
		dat, _ := ioutil.ReadFile(files[0])

		b := bytes.NewReader(dat)

		r, err := zlib.NewReader(b)
		if err != nil {
			panic(err)
		}
		p, _ := ioutil.ReadAll(r)
		fmt.Println(p)

		r.Close()

		return object{data: p, objectType: "blob"}
	}

	return object{}
}

func test() {
	buff := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207,
		47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}
	b := bytes.NewReader(buff)

	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, r)

	r.Close()
}
