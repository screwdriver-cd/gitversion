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

// LastCommit gets the last commit SHA
func LastCommit() (string, error) {
	cmd := execCommand("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("fetching git commit: %v", err)
	}

	trimmed := strings.TrimSpace(string(out))
	return trimmed, nil
}

// LastCommitMessage gets the last commit message
func LastCommitMessage() (string, error) {
	cmd := execCommand("git", "log", "-1", "--pretty=%B")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("fetching git commit message: %v", err)
	}

	trimmed := strings.TrimSpace(string(out))
	return trimmed, nil
}

// Tagged returns true if the specified commit has been tagged
func Tagged() (bool, error) {
	commit, err := LastCommit()
	if err != nil {
		return false, fmt.Errorf("checking current tag: %v", err)
	}
	cmd := execCommand("git", "tag", "--contains", commit)
	t, err := cmd.Output()
	if err != nil {
		return false, nil
	}
	return len(string(t)) > 0, nil
}
