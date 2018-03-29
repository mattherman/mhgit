package objects

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mattherman/mhgit/utils"
)

// Object represents a Git object. It can be of type "blob",
// "commit", or "tree".
type Object struct {
	Data       []byte
	ObjectType string
}

// Index represents the git index
type Index struct {
	Signature  string
	Version    uint32
	EntryCount uint32
	Entries    []IndexEntry
	Checksum   string
}

// IndexEntry represents a file in the git index.
type IndexEntry struct {
	CTimeSec  int32
	CTimeNano int32
	MTimeSec  int32
	MTimeNano int32
	Dev       int32
	Ino       int32
	Mode      int32
	UID       int32
	GID       int32
	FileSize  int32
	Hash      string
	Flags     [2]byte
	Path      string
}

// HashFile will compute the SHA1 hash of a file. If write is true, the
// resulting object will be written to file.
func HashFile(filename string, write bool) (string, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.New("The file was not found")
	}

	obj := Object{Data: fileContent, ObjectType: "blob"}

	return HashObject(obj, write)
}

// HashObject will compute the SHA1 hash of the object and its headers.
// If write is true, the object will be written to file with zlib compression.
func HashObject(objectToHash Object, write bool) (string, error) {
	header := []byte(objectToHash.ObjectType + " " + strconv.Itoa(len(objectToHash.Data)) + "\000")

	fullData := append(header, objectToHash.Data...)

	sha1 := utils.ComputeSha1(fullData)

	if write {
		objectPath := filepath.Join("./.git/objects", sha1[:2])

		os.Mkdir(objectPath, 0700)
		fileName := filepath.Join(objectPath, sha1[2:])

		err := writeCompressedFile(fileName, fullData)

		if err != nil {
			return sha1, errors.New("Object hash calculated, but unable to write to database")
		}
	}

	return sha1, nil
}

// FindObject will attempt to find an object using the given prefix and
// return the filepath to the object. The prefix must at least three
// characters and must be long enough to be unique among all other objects.
func FindObject(hash string) (string, error) {
	if len(hash) < 3 {
		return "", errors.New("Prefix provided must be at least three characters")
	}

	hashDirectory := "./.git/objects/" + hash[:2]
	hashRemainder := hash[2:]

	files, _ := filepath.Glob(hashDirectory + "/" + hashRemainder + "*")

	if len(files) == 0 {
		return "", errors.New("Object " + hash + " not found.")
	} else if len(files) > 1 {
		return "", errors.New("Found multiple matches for " + hash + ".")
	} else {
		return files[0], nil
	}
}

// ReadObject will attempt to find an object using the given prefix and
// return an Object containing the object type and data.
// The prefix must at least three characters and must be long enough to
// be unique among all other objects.
func ReadObject(hash string) (Object, error) {
	objectPath, err := FindObject(hash)

	if err != nil {
		return Object{}, err
	}

	content, err := readCompressedFile(objectPath)

	if err != nil {
		panic(err)
	}

	nullIndex := bytes.IndexByte(content, 0)

	if nullIndex == -1 {
		panic("Object format unreadable.")
	}

	header := string(content[:nullIndex-1])
	headerParts := strings.Split(header, " ")

	if len(headerParts) < 2 {
		panic("Unable to parse object header.")
	}

	return Object{Data: content[nullIndex:], ObjectType: headerParts[0]}, nil
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
