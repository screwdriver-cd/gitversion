package git

import (
	"bufio"
	"fmt"
	"os/exec"
)

var execCommand = exec.Command

// Tags returns the list of git tags as a string slice
func Tags() ([]string, error) {
	cmd := execCommand("git", "tag")
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("fetching git tags: %v", err)
	}

	var lines []string
	scanner := bufio.NewScanner(outReader)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
	}()

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("starting git command: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return nil, fmt.Errorf("waiting for git command: %v", err)
	}

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
