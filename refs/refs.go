package refs

import (
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
		return splitHeadString[2], nil
	}

	return "", nil
}
