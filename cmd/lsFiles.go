package cmd

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/mattherman/mhgit/objects"
	"github.com/mattherman/mhgit/utils"
	"github.com/spf13/cobra"
)

// lsFilesCmd represents the lsFiles command
var lsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "Show information about files in the index and the working tree",
	Run: func(cmd *cobra.Command, args []string) {
		ReadIndex(showStaged)
	},
}

var showStaged bool

func init() {
	rootCmd.AddCommand(lsFilesCmd)
	lsFilesCmd.Flags().BoolVarP(&showStaged, "stage", "s", false, "Show staged contents' mode bits, object name, and stage number in the output.")
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

	headerBytes := indexBytes[0:12]
	checksumBytes := indexBytes[(indexSize - 20):]

	index.Signature = string(headerBytes[0:4])
	index.Version = binary.BigEndian.Uint32(headerBytes[4:8])
	index.EntryCount = binary.BigEndian.Uint32(headerBytes[8:12])
	index.Checksum = hex.EncodeToString(checksumBytes)

	digest := utils.ComputeSha1(indexBytes[:(indexSize - 20)])
	if digest != index.Checksum {
		panic("Index content did not match the checksum")
	}

	index.Entries = make([]objects.IndexEntry, index.EntryCount)
	if index.EntryCount > 0 {
		entryListBytes := indexBytes[12:(indexSize - 20)]
		fmt.Printf("\nindex.EntryCount: %d", index.EntryCount)
		fmt.Printf("\nlen(entryListBytes): %d", len(entryListBytes))

		entryIndex := 0
		for i := 0; i < int(index.EntryCount); i++ {
			entryBytes := entryListBytes[entryIndex:(entryIndex + 62)]
			remainingBytes := entryListBytes[(entryIndex + 62):]
			nullIndex := bytes.IndexByte(remainingBytes, 0)
			pathBytes := remainingBytes[:nullIndex]
			fullEntryBytes := append(entryBytes, pathBytes...)
			index.Entries[i] = readIndexEntry(fullEntryBytes)
			entryIndex += len(fullEntryBytes) + (len(pathBytes) % 8)

			fmt.Printf("\nlen(fullEntryBytes) + (len(pathBytes) %% 8): %d", len(fullEntryBytes)+(len(pathBytes)%8))
			fmt.Printf("\n%+v\n", index.Entries[i])
		}
	}

	return index
}

// Represents an IndexEntry type with only fixed size
// properties for easier parsing. In this type hash property is
// represented as a 20-byte slice instead of the hex string
// and path is ommitted since it is variable length.
type fixedSizeIndexEntry struct {
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
	Hash      [20]byte
	Flags     [2]byte
}

func readIndexEntry(entryBytes []byte) objects.IndexEntry {
	fixedSizeEntry := fixedSizeIndexEntry{}
	entryWithoutPath := entryBytes[:62]

	buf := bytes.NewReader(entryWithoutPath)
	err := binary.Read(buf, binary.BigEndian, &fixedSizeEntry)
	if err != nil {
		panic("Index entry contents malformed")
	}

	path := string(entryBytes[62:])

	return objects.IndexEntry{
		CTimeSec:  fixedSizeEntry.CTimeSec,
		CTimeNano: fixedSizeEntry.CTimeNano,
		MTimeSec:  fixedSizeEntry.MTimeSec,
		MTimeNano: fixedSizeEntry.MTimeNano,
		Dev:       fixedSizeEntry.Dev,
		Ino:       fixedSizeEntry.Ino,
		Mode:      fixedSizeEntry.Mode,
		UID:       fixedSizeEntry.UID,
		GID:       fixedSizeEntry.GID,
		FileSize:  fixedSizeEntry.FileSize,
		Hash:      utils.ComputeSha1(fixedSizeEntry.Hash[:]),
		Path:      path,
	}
}
