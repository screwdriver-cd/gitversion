package main

import (
	"fmt"
	"log"
	"os"

	"github.com/screwdriver-cd/gitversion/bumper"
	"github.com/urfave/cli/v2"
)

// These variables get set by the build script via the LDFLAGS
// Detail about these variables are here: https://goreleaser.com/#builds
var (
	VERSION = "dev"
	COMMIT  = "none"
	DATE    = "unknown"
)

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

	bumpWithFieldAction := func(field bumper.Field) cli.ActionFunc {
		return func(context *cli.Context) error {
			b := bumper.NewBumper()
			return b.Bump(
				bumper.WithPrefix(prefix),
				bumper.WithField(field),
				bumper.WithMerged(merged),
				bumper.WithDryRun(dryrun),
			)
		}
	}

	var latestAction cli.ActionFunc = func(context *cli.Context) error {
		b := bumper.NewBumper()
		v, err := b.LatestVersion(prefix, merged)
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
					Name:   "prerelease",
					Usage:  "bump the prerelease version",
					Action: bumpWithFieldAction(bumper.FieldPrerelease),
				},
				{
					Name:   "patch",
					Usage:  "bump the patch version",
					Action: bumpWithFieldAction(bumper.FieldPatch),
				},
				{
					Name:   "minor",
					Usage:  "bump the minor version",
					Action: bumpWithFieldAction(bumper.FieldMinor),
				},
				{
					Name:   "major",
					Usage:  "bump the major version",
					Action: bumpWithFieldAction(bumper.FieldMajor),
				},
				{
					Name:   "auto",
					Usage:  "bump the version specified in the last commit",
					Action: bumpWithFieldAction(bumper.FieldAuto),
				},
			},
		},
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "output the latest tagged version",
			Action:  latestAction,
		},
	}

	app.Action = latestAction

	// if Run receives an error, the error message is already printed out to
	// stderr, but we should exit with an error code
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
