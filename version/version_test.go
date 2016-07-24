package version

import (
	"fmt"
	"sort"
	"testing"
)

func TestFromString(t *testing.T) {
	var tests = []struct {
		input string
		want  Version
	}{
		{"1.2.3", Version{1, 2, 3}},
		{"3.2.1", Version{3, 2, 1}},
		{"a.b.c", Version{0, 0, 0}},
	}
	for _, test := range tests {
		if v, _ := FromString(test.input); v != test.want {
			t.Errorf("FromString(%q) = %v, want %v", test.input, v, test.want)
		}
	}
}

func TestBadString(t *testing.T) {
	v, err := FromString("a.b.c")
	if err == nil {
		msg := fmt.Sprintf("err should not be nil. v = %v", v)
		t.Error(msg)
	}
}

func TestToString(t *testing.T) {
	want := "1.2.3"
	v, _ := FromString(want)
	got := fmt.Sprintf("%v", v)
	if got != want {
		t.Errorf(`FromString(%q).String() == %q`, want, got)
	}
}

func TestVersionListSort(t *testing.T) {
	var versions = List{
		{2, 1, 3},
		{1, 2, 3},
		{2, 2, 3},
		{3, 1, 2},
		{1, 2, 3},
	}
	var want = List{
		{1, 2, 3},
		{1, 2, 3},
		{2, 1, 3},
		{2, 2, 3},
		{3, 1, 2},
	}

	sort.Sort(versions)
	if len(versions) != len(want) {
		t.Errorf(`error sorting Versions. %v != %v`, versions, want)
	}
	for i, v := range versions {
		if want[i] != v {
			t.Errorf("Version %v != %v", want[i], v)
		}
	}
}
