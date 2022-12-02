package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/screwdriver-cd/gitversion/git"
	"github.com/screwdriver-cd/gitversion/version"
	"github.com/urfave/cli/v2"
)

const (
	// Auto will determine the field based on the last commit message
	Auto = "auto"
	// Major is for specifying the Major field X.0.0
	Major = "major"
	// Minor is for specifying the Minor field 0.X.0
	Minor = "minor"
	// Patch is for specifying the Patch field 0.0.X
	Patch = "patch"
	// PreRelease is for specifying the PreRelease field 0.0.0-X
	PreRelease = "prerelease"
	// MatchField is the index to get the first capture group, aka the field (major, minor, etc.)
	MatchField = 1
)

// These variables get set by the build script via the LDFLAGS
// Detail about these variables are here: https://goreleaser.com/#builds
var (
	VERSION = "dev"
	COMMIT  = "none"
	DATE    = "unknown"
)

var gitTags = git.Tags
var gitTag = git.Tag
var gitTagged = git.Tagged
var gitCommit = git.LastCommit
var gitMessage = git.LastCommitMessage
var errNoVersionTags = errors.New("no valid version tags found")

// Bump increments the specified field of the latest version
func Bump(prefix string, field string, merged, dryrun bool) error {
	v, err := latestVersion(prefix, merged)
	if err != nil {
		if err == errNoVersionTags {
			s := err.Error()
			s = fmt.Sprintf("%s%s", strings.ToUpper(string(s[0])), s[1:])
			log.Printf("WARNING: %v. Using %v", s, v)
		} else {
			return fmt.Errorf("getting latest version %v: %w", v, err)
		}
	}

	log.Printf("Bumping %v for version %v", field, v)
	if field == Auto {
		// If this commit already has a tag, patch
		if tag, _ := gitTagged(); tag {
			field = Patch
		} else {
			// Get commit message and find any reference
			cm, mesErr := gitMessage()
			if mesErr != nil {
				return fmt.Errorf("determing auto patch %w", mesErr)
			}
			re := regexp.MustCompile(`(?i)\[(major|minor|patch|prerelease)( bump)?\]`)
			m := re.FindStringSubmatch(cm)
			if len(m) == 0 {
				field = Patch
			} else {
				field = strings.ToLower(m[MatchField])
			}
		}
	}

	switch field {
	default:
		return errors.New("unknown field type")
	case Major:
		v.Major++
		v.Minor = 0
		v.Patch = 0
	case Minor:
		v.Minor++
		v.Patch = 0
	case Patch:
		v.Patch++
	case PreRelease:
		commit, cerr := gitCommit(true)
		if cerr != nil {
			return fmt.Errorf("getting current commit sha %w", cerr)
		}
		v.PreRelease = commit
	}

	if dryrun {
		log.Print("Dryrun; not git tagging")
	} else if err = gitTag(fmt.Sprintf("%s%s", prefix, v.String())); err != nil {
		return fmt.Errorf("creating new tag %v: %w", v, err)
	}
	// Print out the new tag
	fmt.Printf("%s%s\n", prefix, v)
	return nil
}

func latestVersion(prefix string, merged bool) (v version.Version, err error) {
	versions, err := versions(prefix, merged)
	if err != nil {
		return v, err
	}

	sort.Sort(sort.Reverse(&versions))
	return versions[0], err
}

func versions(prefix string, merged bool) (version.List, error) {
	versions := version.List{}
	tags, err := gitTags(merged)
	if err != nil {
		return nil, fmt.Errorf("fetching git tags: %w", err)
	}

	for _, tag := range tags {
		if len(tag) <= len(prefix) || tag[:len(prefix)] != prefix {
			continue
		}
		tag = tag[len(prefix):]
		v, err := version.FromString(tag)
		if err != nil {
			continue
		}
		versions = append(versions, v)
	}

	if len(versions) == 0 {
		return nil, errNoVersionTags
	}
	return versions, nil
}

func main() {
	var prefix string
	var merged, dryrun bool

	app := cli.NewApp()
	app.Name = "gitversion"
	app.Usage = "manage versions using git tags."
	app.Version = fmt.Sprintf("%v, commit %v, built at %v", VERSION, COMMIT, DATE)

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "prefix",
			Usage:       "set a prefix for the tag name (e.g. v1.0.0)",
			Destination: &prefix,
		},
		&cli.BoolFlag{
			Name:        "merged",
			Usage:       "consider tags merged into this branch",
			Destination: &merged,
		},
	}

	actionLatest := func(c *cli.Context) error {
		v, err := latestVersion(prefix, merged)
		if err != nil {
			log.Printf("Error: %v", err)
			return err
		}
		_, err = fmt.Printf("%s%s\n", prefix, v)
		return err
	}

	app.Commands = []*cli.Command{
		{
			Name:    "bump",
			Aliases: []string{"b"},
			Usage:   "increment the version and create a new git tag",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "dry-run",
					Usage:       "do not add a git tag; only report the tag that would be added",
					Destination: &dryrun,
					Aliases:     []string{"n"},
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:  "prerelease",
					Usage: "bump the prerelease version",
					Action: func(c *cli.Context) error {
						if err := Bump(prefix, PreRelease, merged, dryrun); err != nil {
							log.Printf("Error: %v", err)
							return err
						}
						return nil
					},
				},
				{
					Name:  "patch",
					Usage: "bump the patch version",
					Action: func(c *cli.Context) error {
						if err := Bump(prefix, Patch, merged, dryrun); err != nil {
							log.Printf("Error: %v", err)
							return err
						}
						return nil
					},
				},
				{
					Name:  "minor",
					Usage: "bump the minor version",
					Action: func(c *cli.Context) error {
						if err := Bump(prefix, Minor, merged, dryrun); err != nil {
							log.Printf("Error: %v", err)
							return err
						}
						return nil
					},
				},
				{
					Name:  "major",
					Usage: "bump the major version",
					Action: func(c *cli.Context) error {
						if err := Bump(prefix, Major, merged, dryrun); err != nil {
							log.Printf("Error: %v", err)
							return err
						}
						return nil
					},
				},
				{
					Name:  "auto",
					Usage: "bump the version specified in the last commit",
					Action: func(c *cli.Context) error {
						if err := Bump(prefix, Auto, merged, dryrun); err != nil {
							log.Printf("Error: %v", err)
							return err
						}
						return nil
					},
				},
			},
		},
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "output the latest tagged version",
			Action:  actionLatest,
		},
	}

	app.Action = actionLatest

	// if Run receives an error, the error message is already printed out to
	// stderr, but we should exit with an error code
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
