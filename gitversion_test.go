package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/screwdriver-cd/gitversion/git"
)

func fakeGitTags(bool) ([]string, error) {
	return []string{
		"1.0.1",
		"2.0.1",
		"1.2.1",
		"latest",
		"stable",
		"1.4.3",
		"2.1.2",
	}, nil
}

func TestVersions(t *testing.T) {
	expected := []string{
		"1.0.1",
		"2.0.1",
		"1.2.1",
		"1.4.3",
		"2.1.2",
	}

	gitTags = fakeGitTags
	defer func() { gitTags = git.Tags }()

	v, err := versions("", false)
	require.NoError(t, err)

	assert.Equal(t, len(expected), len(v))

	for i, version := range v {
		assert.Equalf(t, expected[i], version.String(), "Versions()[%d] = %s, want %s", i, version, expected[i])
	}
}

func TestNoVersions(t *testing.T) {
	gitTags = func(bool) ([]string, error) {
		return []string{}, nil
	}
	defer func() { gitTags = git.Tags }()

	v, err := versions("", false)
	require.Error(t, err, "error value for an empty version list should be non-nil")
	assert.Empty(t, v)
}

func TestLatestVersion(t *testing.T) {
	gitTags = fakeGitTags
	defer func() { gitTags = git.Tags }()

	want := "2.1.2"

	latest, err := latestVersion("", false)
	require.NoError(t, err)
	assert.Equal(t, want, latest.String())
}

func TestBumpAutoTagged(t *testing.T) {
	expected := "1.1.2"
	gitTag = func(tag string) error {
		assert.Equal(t, expected, tag)
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	gitTagged = func() (bool, error) {
		return true, nil
	}
	defer func() { gitTagged = git.Tagged }()

	require.NoError(t, Bump("", Auto, false, false))
}

func TestBumpAutoMatch(t *testing.T) {
	expected := "2.0.0"
	gitTag = func(tag string) error {
		assert.Equal(t, expected, tag)
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	gitTagged = func() (bool, error) {
		return false, nil
	}
	defer func() { gitTagged = git.Tagged }()

	gitMessage = func() (string, error) {
		return "[Major] foo", nil
	}
	defer func() { gitMessage = git.LastCommitMessage }()

	require.NoError(t, Bump("", Auto, false, false))
}

func TestBumpAutoMatchAlternate(t *testing.T) {
	expected := "2.0.0"
	gitTag = func(tag string) error {
		assert.Equal(t, expected, tag)
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	gitTagged = func() (bool, error) {
		return false, nil
	}
	defer func() { gitTagged = git.Tagged }()

	gitMessage = func() (string, error) {
		return "[major bump] foo", nil
	}
	defer func() { gitMessage = git.LastCommitMessage }()

	require.NoError(t, Bump("", Auto, false, false))
}

func TestBumpAutoMatchFallback(t *testing.T) {
	expected := "1.1.2"
	gitTag = func(tag string) error {
		assert.Equal(t, expected, tag)
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	gitTagged = func() (bool, error) {
		return false, nil
	}
	defer func() { gitTagged = git.Tagged }()

	gitMessage = func() (string, error) {
		return "foo bar", nil
	}
	defer func() { gitMessage = git.LastCommitMessage }()

	require.NoError(t, Bump("", Auto, false, false))
}

func TestBumpPreRelease(t *testing.T) {
	expected := "1.1.1-9d8ceaa"
	gitTag = func(tag string) error {
		assert.Equal(t, expected, tag)
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitCommit = func(_ bool) (string, error) {
		return "9d8ceaa", nil
	}
	defer func() { gitCommit = git.LastCommit }()

	gitTags = func(bool) ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	require.NoError(t, Bump("", PreRelease, false, false))
}

func TestBumpPatch(t *testing.T) {
	expected := "1.1.2"
	gitTag = func(tag string) error {
		assert.Equal(t, expected, tag)
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	require.NoError(t, Bump("", Patch, false, false))
}

func TestBumpMinor(t *testing.T) {
	expected := "1.2.0"
	gitTag = func(tag string) error {
		assert.Equal(t, expected, tag)
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	require.NoError(t, Bump("", Minor, false, false))
}

func TestBumpMinorDryRun(t *testing.T) {
	gitTag = func(tag string) error {
		assert.Fail(t, "Unexpected call to gitTag")
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	require.NoError(t, Bump("", Minor, false, true))
}

func TestBumpMajor(t *testing.T) {
	expected := "2.0.0"
	gitTag = func(tag string) error {
		assert.Equal(t, expected, tag)
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	require.NoError(t, Bump("", Major, false, false))
}

func TestBumpWithNoVersions(t *testing.T) {
	expected := "0.0.1"
	gitTag = func(tag string) error {
		assert.Equal(t, expected, tag)
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{}, nil
	}
	defer func() { gitTags = git.Tags }()

	require.NoError(t, Bump("", Patch, false, false))
}

func TestBumpWithBadField(t *testing.T) {
	expected := "0.0.1"
	gitTag = func(tag string) error {
		assert.Equal(t, expected, tag)
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{}, nil
	}
	defer func() { gitTags = git.Tags }()

	require.EqualError(t, Bump("", "foobar", false, false), "unknown field type")
}

func TestPrefix(t *testing.T) {
	expected := "2.1.0"
	gitTags = func(bool) ([]string, error) {
		return []string{
			"bigPrefix2.1.0",
			"v2.2.0",
			"2.3.0",
		}, nil
	}
	defer func() { gitTags = git.Tags }()

	latest, err := latestVersion("bigPrefix", false)
	require.NoError(t, err)
	assert.Equal(t, expected, latest.String())
}

func ExampleBump() {
	gitTag = func(tag string) error {
		return nil
	}

	defer func() { gitTag = git.Tag }()

	gitTags = func(bool) ([]string, error) {
		return []string{
			"v2.2.0",
		}, nil
	}
	defer func() { gitTags = git.Tags }()

	Bump("v", Patch, false, false)
	// Output: v2.2.1
}
