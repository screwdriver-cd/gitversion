package main

import (
	"testing"

	"github.com/screwdriver-cd/gitversion/git"
)

func fakeGitTags() ([]string, error) {
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

	v, err := versions("")

	if err != nil {
		t.Errorf("versions() error = %v, should be nil", err)
	}

	if len(expected) != len(v) {
		t.Fatalf("len(Versions()) = %v, want %v", len(v), len(expected))
	}

	for i, version := range v {
		if expected[i] != version.String() {
			t.Errorf("Versions()[%v] = %v, want %v", i, version, expected[i])
		}
	}
}

func TestNoVersions(t *testing.T) {
	gitTags = func() ([]string, error) {
		return []string{}, nil
	}
	defer func() { gitTags = git.Tags }()

	v, err := versions("")
	if err == nil {
		t.Errorf("error value for an empty version list should be non-nil")
	}
	if len(v) != 0 {
		t.Errorf("Expected empty version.List, got: %v", v)
	}
}

func TestLatestVersion(t *testing.T) {
	gitTags = fakeGitTags
	defer func() { gitTags = git.Tags }()

	want := "2.1.2"

	latest, err := latestVersion("")
	if err != nil {
		t.Errorf("Unexpected error from latestVersion(): %v", err)
	}
	if latest.String() != want {
		t.Errorf("latestVersion() = %v, want %v", latest, want)
	}
}

func TestBumpAutoTagged(t *testing.T) {
	expected := "1.1.2"
	gitTag = func(tag string) error {
		if tag != expected {
			t.Errorf("git.Tag() called with %v, want %v", tag, expected)
		}
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func() ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	gitTagged = func() (bool, error) {
		return true, nil
	}
	defer func() { gitTagged = git.Tagged }()

	Bump("", Auto)
}

func TestBumpAutoMatch(t *testing.T) {
	expected := "2.0.0"
	gitTag = func(tag string) error {
		if tag != expected {
			t.Errorf("git.Tag() called with %v, want %v", tag, expected)
		}
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func() ([]string, error) {
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

	Bump("", Auto)
}

func TestBumpAutoMatchAlternate(t *testing.T) {
	expected := "2.0.0"
	gitTag = func(tag string) error {
		if tag != expected {
			t.Errorf("git.Tag() called with %v, want %v", tag, expected)
		}
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func() ([]string, error) {
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

	Bump("", Auto)
}

func TestBumpAutoMatchFallback(t *testing.T) {
	expected := "1.1.2"
	gitTag = func(tag string) error {
		if tag != expected {
			t.Errorf("git.Tag() called with %v, want %v", tag, expected)
		}
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func() ([]string, error) {
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

	Bump("", Auto)
}

func TestBumpPatch(t *testing.T) {
	expected := "1.1.2"
	gitTag = func(tag string) error {
		if tag != expected {
			t.Errorf("git.Tag() called with %v, want %v", tag, expected)
		}
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func() ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	Bump("", Patch)
}

func TestBumpMinor(t *testing.T) {
	expected := "1.2.0"
	gitTag = func(tag string) error {
		if tag != expected {
			t.Errorf("git.Tag() called with %v, want %v", tag, expected)
		}
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func() ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	Bump("", Minor)
}

func TestBumpMajor(t *testing.T) {
	expected := "2.0.0"
	gitTag = func(tag string) error {
		if tag != expected {
			t.Errorf("git.Tag() called with %v, want %v", tag, expected)
		}
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func() ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	Bump("", Major)
}

func TestBumpWithNoVersions(t *testing.T) {
	expected := "0.0.1"
	gitTag = func(tag string) error {
		if tag != expected {
			t.Errorf("git.Tag() called with %v, want %v", tag, expected)
		}
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func() ([]string, error) {
		return []string{}, nil
	}
	defer func() { gitTags = git.Tags }()

	Bump("", Patch)
}

func TestBumpWithBadField(t *testing.T) {
	expected := "0.0.1"
	gitTag = func(tag string) error {
		if tag != expected {
			t.Errorf("git.Tag() called with %v, want %v", tag, expected)
		}
		return nil
	}
	defer func() { gitTag = git.Tag }()

	gitTags = func() ([]string, error) {
		return []string{}, nil
	}
	defer func() { gitTags = git.Tags }()

	err := Bump("", "foobar")
	if err == nil {
		t.Error("expected error from Bump()")
	}
}

func TestPrefix(t *testing.T) {
	expected := "2.1.0"
	gitTags = func() ([]string, error) {
		return []string{
			"bigPrefix2.1.0",
			"v2.2.0",
			"2.3.0",
		}, nil
	}
	defer func() { gitTags = git.Tags }()

	latest, err := latestVersion("bigPrefix")
	if err != nil {
		t.Errorf("unexpected error from latestVersion(): %v", err)
	}

	if latest.String() != expected {
		t.Errorf("latestVersion() = %v, want %v", latest, expected)
	}
}

func ExampleBump() {
	gitTag = func(tag string) error {
		return nil
	}

	defer func() { gitTag = git.Tag }()

	gitTags = func() ([]string, error) {
		return []string{
			"v2.2.0",
		}, nil
	}
	defer func() { gitTags = git.Tags }()

	Bump("v", Patch)
	// Output: v2.2.1
}
