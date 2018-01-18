package command

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/mattherman/mhgit/objects"
	"github.com/mattherman/mhgit/utils"
)

// UpdateIndex will update the current staging area with the current
// state of the file provided
func UpdateIndex(fileName string) {

}

// ReadIndex will show information about files in the
// index and the working tree
func ReadIndex(stage bool) {
	index := readIndexFile()
	fmt.Printf("SIG: %s\nVER: %d\nCOUNT: %d\nCHECK: %s", index.Signature, index.Version, index.EntryCount, index.Checksum)
}

func readIndexFile() objects.Index {
	indexFile := ".git/index"

	if utils.FileDoesNotExist(indexFile) {
		return objects.Index{
			Signature:  "DIRC",
			Version:    0,
			EntryCount: 0,
			Entries:    []objects.IndexEntry{},
		}
	}

	indexBytes, err := ioutil.ReadFile(indexFile)
	if err != nil {
		panic(err)
	}

	indexSize := len(indexBytes)
	index := objects.Index{}

	index.Signature = string(indexBytes[0:4])
	index.Version = binary.BigEndian.Uint32(indexBytes[4:8])
	index.EntryCount = binary.BigEndian.Uint32(indexBytes[8:12])
	index.Checksum = hex.EncodeToString(indexBytes[(indexSize - 20):])

	digest := utils.ComputeSha1(indexBytes[:(indexSize - 20)])
	if digest != index.Checksum {
		panic("Index content did not match the checksum")
	}

	if index.EntryCount > 0 {

	}

	return index
}
