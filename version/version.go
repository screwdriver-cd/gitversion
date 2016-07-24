package version

import (
	"fmt"
	"strconv"
	"strings"
)

// A Version is a version of the form <major>.<minor>.<patch>
type Version struct {
	Major int
	Minor int
	Patch int
}

// Make Builds a Version from major, minor, patch
func Make(major int, minor int, patch int) Version {
	return Version{major, minor, patch}
}

// FromString returns a Version based on a string
func FromString(v string) (ver Version, err error) {
	components := strings.Split(v, ".")
	if len(components) != 3 {
		return ver, fmt.Errorf("Version must contain 3 components: X.Y.Z")
	}

	maj, err := strconv.Atoi(components[0])
	if err != nil {
		return ver, fmt.Errorf("parsing %s as a version: %v", v, err)
	}

	min, err := strconv.Atoi(components[1])
	if err != nil {
		return ver, fmt.Errorf("parsing %s as a version: %v", v, err)
	}

	patch, err := strconv.Atoi(components[2])
	if err != nil {
		return ver, fmt.Errorf("parsing %s as a version: %v", v, err)
	}
	return Version{maj, min, patch}, nil
}

// String formats Version as <major>.<minor>.<patch>
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// List is a slice of Versions that implements sort.Interface
type List []Version

// Len implements sort.Interface.Len()
func (v List) Len() int {
	return len(v)
}

// Less implements sort.Interface.Less()
func (v List) Less(i, j int) bool {
	if v[i].Major != v[j].Major {
		return v[i].Major < v[j].Major
	}
	if v[i].Minor != v[j].Minor {
		return v[i].Minor < v[j].Minor
	}
	if v[i].Patch != v[j].Patch {
		return v[i].Patch < v[j].Patch
	}
	return false
}

// Swap implements sort.Interface.Swap()
func (v List) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
