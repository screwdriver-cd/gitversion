package git

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

type execFunc func(command string, args ...string) *exec.Cmd

func getFakeExecCommand(validator func(string, ...string)) execFunc {
	return func(command string, args ...string) *exec.Cmd {
		validator(command, args...)
		return fakeExecCommand(command, args...)
	}
}

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestTags(t *testing.T) {
	expected := []string{
		"v1.0.1",
		"v2.0.1",
		"v1.2.1",
		"v1.4.3",
	}

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	tags, err := Tags()

	if err != nil {
		t.Errorf("Tags() error = %q, should be nil", err)
	}

	if len(expected) != len(tags) {
		t.Errorf("len(Tags()) = %v, want %v", len(tags), len(expected))
	}

	for i, tag := range tags {
		if expected[i] != tag {
			t.Errorf("Tags()[%v] = %v, want %v", i, tag, expected[i])
		}
	}
}

func TestTag(t *testing.T) {
	execCommand = getFakeExecCommand(func(cmd string, args ...string) {
		if len(args) != 2 {
			t.Errorf("wrong arguments supplied to git tag: %v", args)
		}
		want := "v10.10.10"
		if args[1] != want {
			t.Errorf("Tag received the wrong tag: %q, want %q", args[1], want)
		}
	})
	defer func() { execCommand = exec.Command }()

	Tag("v10.10.10")
}

// This is a fake test for mocking out exec calls.
// See https://golang.org/src/os/exec/exec_test.go and
// https://npf.io/2015/06/testing-exec-command/ for more info
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args[:]
	for i, val := range os.Args { // Should become something lke ["git", "tag"]
		args = os.Args[i:]
		if val == "--" {
			args = args[1:]
			break
		}
	}

	if len(args) >= 2 && args[0] == "git" && args[1] == "tag" {
		if len(args) == 2 {
			tags := []string{
				"v1.0.1",
				"v2.0.1",
				"v1.2.1",
				"v1.4.3",
			}
			for _, tag := range tags {
				fmt.Println(tag)
			}
			return
		}

		if len(args) == 3 {
			return
		}

		os.Exit(255)
	}

	os.Exit(255)
}
