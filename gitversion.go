package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/screwdriver-cd/gitversion/git"
	"github.com/screwdriver-cd/gitversion/version"
	"github.com/urfave/cli"
)

// VERSION gets set by the build script via the LDFLAGS
var VERSION string

// BumpPatch increments the Patch field of the latest version
func BumpPatch(prefix string) error {
	v, err := latestVersion(prefix)
	if err != nil {
		return fmt.Errorf("bumping patch version %v: %v", v, err)
	}

	fmt.Fprintf(os.Stderr, "Bumping patch for version %v\n", v)
	v.Patch++
	if err = gitTag(fmt.Sprintf("%v%v", prefix, v.String())); err != nil {
		return fmt.Errorf("creating new tag %v", v)
	}
	fmt.Fprintf(os.Stdout, "%s%s\n", prefix, v)
	return nil
}

var gitTags = git.Tags
var gitTag = git.Tag

func latestVersion(prefix string) (v version.Version, err error) {
	versions, err := versions(prefix)
	if err != nil {
		return v, err
	}

	sort.Sort(sort.Reverse(&versions))
	return versions[0], err
}

func versions(prefix string) (version.List, error) {
	versions := version.List{}
	tags, err := gitTags()
	if err != nil {
		return nil, fmt.Errorf("fetching git tags: %v", err)
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
		return nil, fmt.Errorf("no valid version tags found")
	}
	return versions, nil
}

func main() {
	var prefix string

	app := cli.NewApp()
	app.Name = "gitversion"
	app.Usage = "manage versions using git tags."
	if VERSION == "" {
		VERSION = "0.0.0"
	}
	app.Version = VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "prefix",
			Usage:       "set a prefix for the tag name (e.g. v1.0.0)",
			Destination: &prefix,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "bump",
			Aliases: []string{"b"},
			Usage:   "increment the version and create a new git tag",
			Subcommands: []cli.Command{
				{
					Name:  "patch",
					Usage: "bump the patch version",
					Action: func(c *cli.Context) error {
						if err := BumpPatch(prefix); err != nil {
							fmt.Fprintf(os.Stderr, "Error: %v", err)
							return err
						}
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
