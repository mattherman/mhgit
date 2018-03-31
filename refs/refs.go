package refs

import (
	"errors"
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

// ListBranches will return all the existing branches
func ListBranches() ([]string, error) {
	files, err := ioutil.ReadDir(".git/refs/heads")

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}
	return filenames, err
}

// CreateBranch will create a new branch if it does not exist
func CreateBranch(branchName string) error {
	branches, err := ListBranches()
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if branch == branchName {
			return errors.New("branch already exists")
		}
	}

	commitHash, err := LatestCommit()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(".git/refs/heads/%s", branchName)
	err = ioutil.WriteFile(filename, []byte(commitHash), 0644)

	return err
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

	return string(bytes), nil
}

// UpdateLatestCommit will update the latest commit of the current branch
// to equal the provided hash.
func UpdateLatestCommit(commitHash string) error {
	branch, err := CurrentBranch()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(".git/refs/heads/%s", branch)
	err = ioutil.WriteFile(filename, []byte(commitHash), 0644)
	if err != nil {
		return err
	}

	return nil
}
