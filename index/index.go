package index

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/mattherman/mhgit/objects"
)

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
	Flags     uint16
	Path      string
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
	Flags     uint16
}

func (e fixedSizeIndexEntry) getPathLength() int {
	return int(e.Flags & 0x0FFF)
}

func (e fixedSizeIndexEntry) toFullEntry(path string) IndexEntry {
	return IndexEntry{
		CTimeSec:  e.CTimeSec,
		CTimeNano: e.CTimeNano,
		MTimeSec:  e.MTimeSec,
		MTimeNano: e.MTimeNano,
		Dev:       e.Dev,
		Ino:       e.Ino,
		Mode:      e.Mode,
		UID:       e.UID,
		GID:       e.GID,
		FileSize:  e.FileSize,
		Hash:      hex.EncodeToString(e.Hash[:]),
		Path:      path,
	}
}

const fixedSizeIndexEntryLength int = 62

// ReadIndex will show information about files in the
// index and the working tree
func ReadIndex() Index {
	indexFile := ".git/index"

	_, err := os.Stat(indexFile)
	if os.IsNotExist(err) {
		return Index{
			Signature:  "DIRC",
			Version:    2,
			EntryCount: 0,
			Entries:    []IndexEntry{},
		}
	}

	indexBytes, err := ioutil.ReadFile(indexFile)
	if err != nil {
		panic(err)
	}

	indexSize := len(indexBytes)
	index := Index{}

	headerBytes := indexBytes[0:12]
	checksumBytes := indexBytes[(indexSize - 20):]

	index.Signature = string(headerBytes[0:4])
	index.Version = binary.BigEndian.Uint32(headerBytes[4:8])
	index.EntryCount = binary.BigEndian.Uint32(headerBytes[8:12])
	index.Checksum = hex.EncodeToString(checksumBytes)

	digest := objects.ComputeSha1(indexBytes[:(indexSize - 20)])
	if digest != index.Checksum {
		panic("Index content did not match the checksum")
	}

	index.Entries = make([]IndexEntry, index.EntryCount)
	if index.EntryCount > 0 {

		entryListBytes := indexBytes[12:(indexSize - 20)]

		entryIndex := 0
		for i := 0; i < int(index.EntryCount); i++ {
			// Convert fixed size portion of the entry to a fixedSizeIndexEntry
			fixedSizeEntryBytes := entryListBytes[entryIndex:(entryIndex + fixedSizeIndexEntryLength)]
			fixedSizeIndexEntry := readIndexEntry(fixedSizeEntryBytes)

			// Get bytes for index entry's path field
			startPathIndex := entryIndex + fixedSizeIndexEntryLength
			pathLength := fixedSizeIndexEntry.getPathLength()
			entryPathBytes := entryListBytes[startPathIndex:(startPathIndex + pathLength)]

			// Convert the fixedSizeIndexEntry + path to a full IndexEntry
			index.Entries[i] = fixedSizeIndexEntry.toFullEntry(string(entryPathBytes))

			// Advance the entry index by the length of the previous entry plus enough
			// null padding to extend the entry to a multiple of 8 bytes
			totalEntryLength := len(fixedSizeEntryBytes) + pathLength
			entryIndex += totalEntryLength + nullPaddingLength(totalEntryLength)
		}
	}

	return index
}

// Determines the number of bytes necessary to extend the given
// byte length to a multiple of 8 while ensuring it is suffixed
// by at least one null byte.
// Ex.
// 		nullPaddingLength(7) => 1
// 		nullPaddingLength(8) => 8
// 		nullPaddingLength(9) => 7
func nullPaddingLength(pathLength int) int {
	return 8 - (pathLength % 8)
}

func readIndexEntry(entryBytes []byte) fixedSizeIndexEntry {
	fixedSizeEntry := fixedSizeIndexEntry{}
	entryWithoutPath := entryBytes[:62]

	buf := bytes.NewReader(entryWithoutPath)
	err := binary.Read(buf, binary.BigEndian, &fixedSizeEntry)
	if err != nil {
		panic("Index entry contents malformed")
	}

	return fixedSizeEntry
}
