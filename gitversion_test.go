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

func TestBumpPatch(t *testing.T) {
	expected := "1.1.2"
	gitTag = func(tag string) error {
		if tag != expected {
			t.Errorf("git.Tag() called with %v, want %v", tag, expected)
		}
		return nil
	}

	gitTags = func() ([]string, error) {
		return []string{"1.1.1", "0.1.1"}, nil
	}
	defer func() { gitTags = git.Tags }()

	BumpPatch("")
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

	latest, err := latestVersion("bigPrefix")
	if err != nil {
		t.Errorf("unexpected error from latestVersion(): %v", err)
	}

	if latest.String() != expected {
		t.Errorf("latestVersion() = %v, want %v", latest, expected)
	}
}

func ExampleBumpPatch() {
	gitTags = func() ([]string, error) {
		return []string{
			"v2.2.0",
		}, nil
	}

	BumpPatch("v")
	// Output: v2.2.1
}
