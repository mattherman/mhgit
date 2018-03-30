package index

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"os"
	"sort"
	"syscall"

	"github.com/mattherman/mhgit/objects"
)

const fixedSizeIndexEntryLength int = 62
const checksumLength int = 20
const indexFile string = ".git/index"

// Index represents the git index
type Index struct {
	Signature  string
	Version    uint32
	EntryCount uint32
	Entries    []Entry
	Checksum   string
}

// Entry represents a file in the git index.
type Entry struct {
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

func (e fixedSizeIndexEntry) toFullEntry(path string) Entry {
	return Entry{
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

func (e Entry) toFixedSizeEntry(path string) fixedSizeIndexEntry {

	var hashArray [20]byte
	hashBytes, _ := hex.DecodeString(e.Hash)
	copy(hashArray[:], hashBytes)

	flags := uint16(len(path)) & 0x0FFF

	return fixedSizeIndexEntry{
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
		Hash:      hashArray,
		Flags:     flags,
	}
}

// NewEntry will create a new index entry based on the filepath given.
// The hash of the file will be included in the entry, but no object
// will be created in the database.
func newEntry(filepath string, hash string) (Entry, error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return Entry{}, err
	}

	var ctimesec int32
	var ctimenano int32
	var mtimesec int32
	var mtimenano int32
	var ino int32
	var dev int32
	var uid int32
	var gid int32
	var mode int32

	statUnix, infoIsAvailable := stat.Sys().(*syscall.Stat_t)
	if infoIsAvailable {
		ctimesec = int32(statUnix.Ctim.Sec)
		ctimenano = int32(statUnix.Ctim.Nsec)
		mtimesec = int32(statUnix.Mtim.Sec)
		mtimenano = int32(statUnix.Mtim.Nsec)
		ino = int32(statUnix.Ino)
		dev = int32(statUnix.Dev)
		uid = int32(statUnix.Uid)
		gid = int32(statUnix.Gid)
		mode = int32(statUnix.Mode)
	}

	return Entry{
		CTimeSec:  ctimesec,
		CTimeNano: ctimenano,
		MTimeSec:  mtimesec,
		MTimeNano: mtimenano,
		Dev:       dev,
		Ino:       ino,
		UID:       uid,
		GID:       gid,
		Mode:      mode,
		FileSize:  int32(stat.Size()),
		Hash:      hash,
		Path:      filepath,
	}, nil
}

// ReadIndex will show information about files in the
// index and the working tree
func ReadIndex() Index {
	_, err := os.Stat(indexFile)
	if os.IsNotExist(err) {
		return Index{
			Signature:  "DIRC",
			Version:    2,
			EntryCount: 0,
			Entries:    []Entry{},
		}
	}

	indexBytes, err := ioutil.ReadFile(indexFile)
	if err != nil {
		panic(err)
	}

	indexSize := len(indexBytes)
	index := Index{}

	headerBytes := indexBytes[0:12]
	checksumBytes := indexBytes[(indexSize - checksumLength):]

	index.Signature = string(headerBytes[0:4])
	index.Version = binary.BigEndian.Uint32(headerBytes[4:8])
	index.EntryCount = binary.BigEndian.Uint32(headerBytes[8:12])
	index.Checksum = hex.EncodeToString(checksumBytes)

	digest := objects.ComputeSha1(indexBytes[:(indexSize - checksumLength)])
	if digest != index.Checksum {
		panic("Index content did not match the checksum")
	}

	index.Entries = make([]Entry, index.EntryCount)
	if index.EntryCount > 0 {

		entryListBytes := indexBytes[12:(indexSize - checksumLength)]

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

// WriteIndex will write the index file with the specified entries
func WriteIndex(entries []Entry) error {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	index := Index{
		Signature:  "DIRC",
		Version:    2,
		EntryCount: uint32(len(entries)),
		Entries:    entries,
	}

	f, err := os.Create(indexFile)
	defer f.Close()
	if err != nil {
		return err
	}

	var header [12]byte
	copy(header[0:4], index.Signature)
	binary.BigEndian.PutUint32(header[4:8], index.Version)
	binary.BigEndian.PutUint32(header[8:12], index.EntryCount)

	var entryBuffer bytes.Buffer
	for _, entry := range index.Entries {
		// Write the fixed size portion of the entry to the buffer, followed by the path
		binary.Write(&entryBuffer, binary.BigEndian, entry.toFixedSizeEntry(entry.Path))
		binary.Write(&entryBuffer, binary.BigEndian, []byte(entry.Path))

		// Add enough null padding to extend the entry to a multiple of 8 bytes with null-termination
		entryLength := fixedSizeIndexEntryLength + len(entry.Path)
		binary.Write(&entryBuffer, binary.BigEndian, make([]byte, nullPaddingLength(entryLength)))
	}

	indexAndEntries := append(header[:], entryBuffer.Bytes()...)
	checksum := objects.ComputeSha1(indexAndEntries)
	checksumBytes, err := hex.DecodeString(checksum)
	if err != nil {
		return err
	}

	fullIndex := append(indexAndEntries, checksumBytes...)
	_, err = f.Write(fullIndex)
	if err != nil {
		return err
	}

	return nil
}

// Exists returns whether or not the provided filepath is present
// in the index
func Exists(filepath string) bool {
	index := ReadIndex()
	for _, entry := range index.Entries {
		if filepath == entry.Path {
			return true
		}
	}
	return false
}

// Add will add the specified file to the index if it exists
// in the working directory
func Add(filepath string) error {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return err
	}

	hash, err := objects.HashFile(filepath, true)
	if err != nil {
		return err
	}

	entry, err := newEntry(filepath, hash)
	if err != nil {
		return err
	}

	// TODO add to index instead of overwriting it
	err = WriteIndex([]Entry{entry})
	if err != nil {
		return err
	}

	return nil
}
