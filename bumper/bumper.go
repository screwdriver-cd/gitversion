package bumper

import (
	"errors"
	"fmt"
	"github.com/screwdriver-cd/gitversion/git"
	"github.com/screwdriver-cd/gitversion/version"
	"log"
	"regexp"
	"sort"
	"strings"
)

//go:generate go run github.com/golang/mock/mockgen -source $GOFILE -destination mock_$GOFILE -package $GOPACKAGE

type (
	bumpOptions struct {
		prefix string
		field  Field
		merged bool
		dryrun bool
	}
	BumpOption func(*bumpOptions)

	Bumper interface {
		Bump(...BumpOption) error
		LatestVersion(prefix string, merged bool) (v version.Version, err error)
		Versions(prefix string, merged bool) (version.List, error)
	}
	DefaultBumper struct {
		Git git.Git
	}
)

var (
	_ Bumper = &DefaultBumper{}

	defaultBumpOptions = []BumpOption{
		WithField(FieldAuto),
	}
)

const (
	// MatchField is the index to get the first capture group, aka the field (major, minor, etc.)
	MatchField = 1
)

func newBumpOptions(options ...BumpOption) *bumpOptions {
	ret := &bumpOptions{}
	for _, opt := range defaultBumpOptions {
		opt(ret)
	}
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

func WithPrefix(prefix string) BumpOption {
	return func(options *bumpOptions) {
		options.prefix = prefix
	}
}

func WithField(field Field) BumpOption {
	return func(options *bumpOptions) {
		options.field = field
	}
}

func WithMerged(merged bool) BumpOption {
	return func(options *bumpOptions) {
		options.merged = merged
	}
}

func WithDryRun(dryrun bool) BumpOption {
	return func(options *bumpOptions) {
		options.dryrun = dryrun
	}
}

var (
	_ Bumper = &DefaultBumper{}

	errNoVersionTags = errors.New("no valid version tags found")
)

func (d *DefaultBumper) Bump(options ...BumpOption) error {
	opts := newBumpOptions(options...)
	field := opts.field

	v, err := d.LatestVersion(opts.prefix, opts.merged)
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
	if field == FieldAuto {
		// If this commit already has a tag, patch
		if tag, _ := d.Git.Tagged(); tag {
			field = FieldPatch
		} else {
			// Get commit message and find any reference
			cm, mesErr := d.Git.LastCommitMessage()
			if mesErr != nil {
				return fmt.Errorf("determing auto patch %w", mesErr)
			}
			re := regexp.MustCompile(`(?i)\[(major|minor|patch|prerelease)( bump)?\]`)
			m := re.FindStringSubmatch(cm)
			if len(m) == 0 {
				field = FieldPatch
			} else {
				if field, err = ParseField(strings.ToLower(m[MatchField])); err != nil {
					return err
				}
			}
		}
	}

	switch field {
	default:
		return errors.New("unknown field type")
	case FieldMajor:
		v.Major++
		v.Minor = 0
		v.Patch = 0
	case FieldMinor:
		v.Minor++
		v.Patch = 0
	case FieldPatch:
		v.Patch++
	case FieldPrerelease:
		commit, cerr := d.Git.LastCommit(true)
		if cerr != nil {
			return fmt.Errorf("getting current commit sha %w", cerr)
		}
		v.PreRelease = commit
	}

	newTag := fmt.Sprintf("%s%s", opts.prefix, v)
	if opts.dryrun {
		log.Print("Dryrun; not git tagging")
	} else if err = d.Git.Tag(newTag); err != nil {
		return fmt.Errorf("creating new tag %v: %w", v, err)
	}

	// Print out the new tag
	_, err = fmt.Println(newTag)
	return err
}

func (d *DefaultBumper) LatestVersion(prefix string, merged bool) (v version.Version, err error) {
	versions, err := d.Versions(prefix, merged)
	if err != nil {
		return v, err
	}

	sort.Sort(sort.Reverse(&versions))
	return versions[0], err
}

func (d *DefaultBumper) Versions(prefix string, merged bool) (version.List, error) {
	versions := version.List{}
	tags, err := d.Git.Tags(merged)
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
