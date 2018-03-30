package refs

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// CurrentBranch returns the name of the branch currently
// pointed to by HEAD or empty string if the ref is not a branch
func CurrentBranch() (string, error) {
	bytes, err := ioutil.ReadFile(".git/HEAD")
	if err != nil {
		return "", err
	}

	headString := string(bytes)
	match := strings.HasPrefix(headString, "ref: refs/heads/")

	if match {
		splitHeadString := strings.Split(headString, "/")
		return strings.Trim(splitHeadString[2], " \n"), nil
	}

	return "", nil
}

// LatestCommit will return the latest commit hash of the current branch
func LatestCommit() (string, error) {
	branch, err := CurrentBranch()
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf(".git/refs/heads/%s", branch)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.Trim(string(bytes), " \n"), nil
}

// UpdateLatestCommit will update the latest commit of the current branch
// to equal the provided hash.
func UpdateLatestCommit(commitHash string) error {
	branch, err := CurrentBranch()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(".git/refs/heads/%s", branch)
	err = ioutil.WriteFile(filename, []byte(commitHash+"\n"), 0644)
	if err != nil {
		return err
	}

	return nil
}
