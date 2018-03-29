package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/mattherman/mhgit/objects"

	"github.com/mattherman/mhgit/index"
	"github.com/spf13/cobra"
)

// updateIndexCmd represents the updateIndex command
var updateIndexCmd = &cobra.Command{
	Use:   "update-index [file]",
	Short: "Register file contents in the working tree to the index.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		addToIndex(filepath)
	},
}

func init() {
	rootCmd.AddCommand(updateIndexCmd)
}

func addToIndex(filepath string) {
	entry, err := createIndexEntryForFile(filepath)
	if err != nil {
		fmt.Printf("Failed to create the index entry: %v\n", err)
	}

	err = index.WriteIndex([]index.IndexEntry{entry})
	if err != nil {
		fmt.Printf("Failed to write to the index: %v\n", err)
	}
}

func createIndexEntryForFile(filepath string) (index.IndexEntry, error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return index.IndexEntry{}, err
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

	hash, err := objects.HashFile(filepath, true)
	if err != nil {
		return index.IndexEntry{}, err
	}

	return index.IndexEntry{
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
