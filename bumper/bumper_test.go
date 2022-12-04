package bumper

import (
	"github.com/golang/mock/gomock"
	"github.com/screwdriver-cd/gitversion/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	MockGitOption func(*git.MockGit)
)

func mockGitForTest(ctrl *gomock.Controller, options ...MockGitOption) *git.MockGit {
	ret := git.NewMockGit(ctrl)
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

func bumperForTest(ctrl *gomock.Controller, options ...MockGitOption) Bumper {
	return &DefaultBumper{
		Git: mockGitForTest(ctrl, options...),
	}
}

func withGitTags(tags ...string) MockGitOption {
	return func(mockGit *git.MockGit) {
		mockGit.EXPECT().
			Tags(gomock.Any()).
			Return(tags, nil)
	}
}
func withFakeGitTags() MockGitOption {
	return withGitTags(
		"1.0.1",
		"2.0.1",
		"1.2.1",
		"latest",
		"stable",
		"1.4.3",
		"2.1.2",
	)
}

func withEmptyGitTags() MockGitOption {
	return func(mockGit *git.MockGit) {
		mockGit.EXPECT().
			Tags(gomock.Any()).
			Return(nil, nil)
	}
}

func withExpectedTag(tag string) MockGitOption {
	return func(mockGit *git.MockGit) {
		mockGit.EXPECT().
			Tag(gomock.Eq(tag))
	}
}

func withTagged(tagged bool) MockGitOption {
	return func(mockGit *git.MockGit) {
		mockGit.EXPECT().
			Tagged().
			Return(tagged, nil)
	}
}

func withLastCommitMessage(message string) MockGitOption {
	return func(mockGit *git.MockGit) {
		mockGit.EXPECT().
			LastCommitMessage().
			Return(message, nil)
	}
}

func withLastCommit(commit string) MockGitOption {
	return func(mockGit *git.MockGit) {
		mockGit.EXPECT().
			LastCommit(gomock.Any()).
			Return(commit, nil)
	}
}

func TestVersions(t *testing.T) {
	ctrl := gomock.NewController(t)

	expected := []string{
		"1.0.1",
		"2.0.1",
		"1.2.1",
		"1.4.3",
		"2.1.2",
	}
	b := bumperForTest(ctrl, withFakeGitTags())

	v, err := b.Versions("", false)
	require.NoError(t, err)

	assert.Equal(t, len(expected), len(v))

	for i, version := range v {
		assert.Equalf(t, expected[i], version.String(), "Versions()[%d] = %s, want %s", i, version, expected[i])
	}
}

func TestNoVersions(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(ctrl, withEmptyGitTags())

	v, err := b.Versions("", false)
	require.Error(t, err, "error value for an empty version list should be non-nil")
	assert.Empty(t, v)
}

func TestLatestVersion(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(ctrl, withFakeGitTags())

	want := "2.1.2"

	latest, err := b.LatestVersion("", false)
	require.NoError(t, err)
	assert.Equal(t, want, latest.String())
}

func TestBumpAutoTagged(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withExpectedTag("1.1.2"),
		withGitTags("1.1.1", "0.1.1"),
		withTagged(true),
	)

	require.NoError(t, b.Bump(WithField(FieldAuto)))
}

func TestBumpAutoMatch(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withExpectedTag("2.0.0"),
		withGitTags("1.1.1", "0.1.1"),
		withTagged(false),
		withLastCommitMessage("[Major] foo"),
	)

	require.NoError(t, b.Bump(WithField(FieldAuto)))
}

func TestBumpAutoMatchAlternate(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withExpectedTag("2.0.0"),
		withGitTags("1.1.1", "0.1.1"),
		withTagged(false),
		withLastCommitMessage("[major bump] foo"),
	)

	require.NoError(t, b.Bump(WithField(FieldAuto)))
}

func TestBumpAutoMatchFallback(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withExpectedTag("1.1.2"),
		withGitTags("1.1.1", "0.1.1"),
		withTagged(false),
		withLastCommitMessage("foo bar"),
	)

	require.NoError(t, b.Bump(WithField(FieldAuto)))
}

func TestBumpPreRelease(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withExpectedTag("1.1.1-9d8ceaa"),
		withLastCommit("9d8ceaa"),
		withGitTags("1.1.1", "0.1.1"),
	)

	require.NoError(t, b.Bump(WithField(FieldPrerelease)))
}

func TestBumpPatch(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withExpectedTag("1.1.2"),
		withGitTags("1.1.1", "0.1.1"),
	)

	require.NoError(t, b.Bump(WithField(FieldPatch)))
}

func TestBumpMinor(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withExpectedTag("1.2.0"),
		withGitTags("1.1.1", "0.1.1"),
	)

	require.NoError(t, b.Bump(WithField(FieldMinor)))
}

func TestBumpMinorDryRun(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withGitTags("1.1.1", "0.1.1"),
	)

	require.NoError(t, b.Bump(WithField(FieldMinor), WithDryRun(true)))
}

func TestBumpMajor(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withExpectedTag("2.0.0"),
		withGitTags("1.1.1", "0.1.1"),
	)

	require.NoError(t, b.Bump(WithField(FieldMajor)))
}

func TestBumpWithNoVersions(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withExpectedTag("0.0.1"),
		withEmptyGitTags(),
	)

	require.NoError(t, b.Bump(WithField(FieldPatch)))
}

func TestBumpWithBadField(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withEmptyGitTags(),
	)

	assert.EqualError(t, b.Bump(WithField("foobar")), "unknown field type")
}

func TestPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)

	b := bumperForTest(
		ctrl,
		withGitTags(
			"bigPrefix2.1.0",
			"v2.2.0",
			"2.3.0",
		),
	)

	latest, err := b.LatestVersion("bigPrefix", false)
	require.NoError(t, err)
	assert.Equal(t, "2.1.0", latest.String())
}

func ExampleBump() {
	t := &testing.T{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := bumperForTest(
		ctrl,
		withExpectedTag("v2.2.1"),
		withGitTags("v2.2.0"),
	)

	require.NoError(t, b.Bump(WithPrefix("v"), WithField(FieldPatch)))
	// Output: v2.2.1
}
