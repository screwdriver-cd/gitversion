package git

import (
	"fmt"
	"os/exec"
	"strings"
)

//go:generate go run github.com/golang/mock/mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

type (
	CmdRunner interface {
		Run(*exec.Cmd) error
		Output(*exec.Cmd) ([]byte, error)
	}
	DefaultCmdRunner struct{}

	Git interface {
		LastCommit(short bool) (string, error)
		LastCommitMessage() (string, error)
		Tag(tag string) error
		Tags(merged bool) ([]string, error)
		Tagged() (bool, error)
	}
	DefaultGit struct {
		CmdRunner CmdRunner
	}
)

var (
	_ Git       = &DefaultGit{}
	_ CmdRunner = &DefaultCmdRunner{}
)

// Tags returns the list of git tags as a string slice
func (g *DefaultGit) Tags(merged bool) ([]string, error) {
	args := []string{"tag"}
	if merged {
		args = append(args, "--merged")
	}
	cmd := exec.Command("git", args...)
	out, err := g.CmdRunner.Output(cmd)
	if err != nil {
		return nil, fmt.Errorf("fetching git tags: %w", err)
	}

	trimmed := strings.TrimSpace(string(out))
	lines := strings.Split(trimmed, "\n")
	return lines, nil
}

// Tag calls git to create a new tag from a string
func (g *DefaultGit) Tag(tag string) error {
	cmd := exec.Command("git", "tag", tag)
	_, err := g.CmdRunner.Output(cmd)
	if err != nil {
		return fmt.Errorf("tagging the commit in git: %w", err)
	}

	return nil
}

// LastCommit gets the last commit SHA
func (g *DefaultGit) LastCommit(short bool) (string, error) {
	var cmd *exec.Cmd

	if short {
		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	} else {
		cmd = exec.Command("git", "rev-parse", "HEAD")
	}
	out, err := g.CmdRunner.Output(cmd)
	if err != nil {
		return "", fmt.Errorf("fetching git commit: %w", err)
	}

	trimmed := strings.TrimSpace(string(out))
	return trimmed, nil
}

// LastCommitMessage gets the last commit message
func (g *DefaultGit) LastCommitMessage() (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	out, err := g.CmdRunner.Output(cmd)
	if err != nil {
		return "", fmt.Errorf("fetching git commit message: %w", err)
	}

	trimmed := strings.TrimSpace(string(out))
	return trimmed, nil
}

// Tagged returns true if the specified commit has been tagged
func (g *DefaultGit) Tagged() (bool, error) {
	commit, err := g.LastCommit(false)
	if err != nil {
		return false, fmt.Errorf("checking current tag: %w", err)
	}
	cmd := exec.Command("git", "tag", "--contains", commit)
	t, err := g.CmdRunner.Output(cmd)
	if err != nil {
		return false, nil
	}
	return len(string(t)) > 0, nil
}

func (d *DefaultCmdRunner) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}

func (d *DefaultCmdRunner) Output(cmd *exec.Cmd) ([]byte, error) {
	return cmd.Output()
}
