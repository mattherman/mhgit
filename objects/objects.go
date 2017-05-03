package objects

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

// Object represents a Git object. It can be of type "blob",
// "commit", or "tree".
type Object struct {
	Data       []byte
	ObjectType string
}

// HashObject will compute the SHA1 hash of the object and its headers.
// If write is true, the object will be written to file with zlib compression.
func HashObject(objectToHash Object, write bool) string {
	header := []byte(objectToHash.ObjectType + " " + strconv.Itoa(len(objectToHash.Data)) + "\000")

	fullData := append(header, objectToHash.Data...)

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

// ReadObject will attempt to find an object using the given prefix.
// The prefix must at least three characters and must be long enough to
// be unique among all other objects.
func ReadObject(hash string) (Object, error) {
	if len(hash) < 3 {
		return Object{}, errors.New("Prefix provided must be at least three characters")
	}

	hashDirectory := "./.git/objects/" + hash[:2]
	hashRemainder := hash[2:]

	files, _ := filepath.Glob(hashDirectory + "/" + hashRemainder + "*")

	if len(files) == 0 {
		return Object{}, errors.New("Object " + hash + " not found.")
	} else if len(files) > 1 {
		return Object{}, errors.New("Found multiple matches for " + hash + ".")
	} else {
		content, err := readCompressedFile(files[0])

		if err != nil {
			panic(err)
		}

		return Object{Data: content, ObjectType: "blob"}, nil
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
