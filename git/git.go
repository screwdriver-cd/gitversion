package git

import (
	"fmt"
	"os/exec"
	"strings"
)

var execCommand = exec.Command

// Tags returns the list of git tags as a string slice
func Tags() ([]string, error) {
	cmd := execCommand("git", "tag")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("fetching git tags: %v", err)
	}

	trimmed := strings.TrimSpace(string(out))
	lines := strings.Split(trimmed, "\n")
	return lines, nil
}

// Tag calls git to create a new tag from a string
func Tag(tag string) error {
	cmd := execCommand("git", "tag", tag)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("tagging the commit in git: %v", err)
	}

	return nil
}
