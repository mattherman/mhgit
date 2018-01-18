package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
)

// FileDoesNotExist returns a bool with the
// value true if the file at the specified path
// does not exist, else returns false.
func FileDoesNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// ComputeSha1 returns the SHA-1 hash of the
// provided byte array as a string.
func ComputeSha1(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}
