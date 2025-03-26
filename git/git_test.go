package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/screwdriver-cd/gitversion/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	CmdRunnerOption func(*MockCmdRunner)
)

var (
	fakeTags = []string{
		"v1.0.1",
		"v2.0.1",
		"v1.2.1",
		"v1.4.3",
	}
	fakeTagsOutput = strings.Join(fakeTags, "\n") + "\n"
	fakeHead       = "9d8ceaaa28f0563e52e1edf3eaae72c814aa1102"
	fakeHeadOutput = fakeHead + "\n"
)

func mockRunnerForTest(ctrl *gomock.Controller, options ...CmdRunnerOption) *MockCmdRunner {
	ret := NewMockCmdRunner(ctrl)
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

func gitForTest(ctrl *gomock.Controller, options ...CmdRunnerOption) Git {
	ret := &DefaultGit{
		CmdRunner: mockRunnerForTest(ctrl, options...),
	}
	return ret
}

func gitCmdMatcher(args ...string) gomock.Matcher {
	return testutil.NewMatcherFunc(fmt.Sprint(args), func(x interface{}) bool {
		cmd := x.(*exec.Cmd)
		if filepath.Base(cmd.Args[0]) != "git" {
			return false
		}
		if len(args)+1 != len(cmd.Args) {
			return false
		}
		for i := range args {
			if cmd.Args[i+1] != args[i] {
				return false
			}
		}
		return true
	})
}

func withGitTagOutput(output string, args ...string) CmdRunnerOption {
	return func(runner *MockCmdRunner) {
		runner.EXPECT().
			Output(gitCmdMatcher(args...)).
			Return([]byte(output), nil)
	}
}

func TestTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := []string{
		"v1.0.1",
		"v2.0.1",
		"v1.2.1",
		"v1.4.3",
	}
	g := gitForTest(ctrl, withGitTagOutput(fakeTagsOutput, "tag"))

	tags, err := g.Tags(false)
	require.NoError(t, err)

	require.Equal(t, len(expected), len(tags))

	for i, tag := range tags {
		assert.Equalf(t, expected[i], tag, "Tags()[%v] = %v, want %v", i, tag, expected[i])
	}
}

func TestLastCommitLong(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := "9d8ceaaa28f0563e52e1edf3eaae72c814aa1102"
	g := gitForTest(ctrl, withGitTagOutput(fakeHeadOutput, "rev-parse", "HEAD"))

	commit, err := g.LastCommit(false)
	require.NoError(t, err)

	assert.Equal(t, expected, commit)
}

func TestLastCommitShort(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := "9d8ceaa"
	g := gitForTest(ctrl, withGitTagOutput(expected+"\n", "rev-parse", "--short", "HEAD"))

	commit, err := g.LastCommit(true)
	require.NoError(t, err)
	if err != nil {
		t.Errorf("LastCommit() error = %q, should be nil", err)
	}

	assert.Equal(t, expected, commit)
}

func TestLastMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := "minor: this should be bumping minor"
	g := gitForTest(ctrl, withGitTagOutput(expected+"\n", "log", "-1", "--pretty=%B"))

	commit, err := g.LastCommitMessage()
	require.NoError(t, err)

	assert.Equal(t, expected, commit)
}

func TestTagged(t *testing.T) {
	ctrl := gomock.NewController(t)
	g := gitForTest(ctrl,
		withGitTagOutput("true\n", "tag", "--contains", fakeHead),
		withGitTagOutput(fakeHeadOutput, "rev-parse", "HEAD"),
	)

	commit, err := g.Tagged()
	require.NoError(t, err)

	require.True(t, commit)
}

func TestTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	expected := "v10.10.10"
	g := gitForTest(ctrl,
		withGitTagOutput("", "tag", expected),
	)

	require.NoError(t, g.Tag(expected))
}
