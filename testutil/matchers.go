package testutil

import "github.com/golang/mock/gomock"

type (
	matcherFunc        func(x interface{}) bool
	matcherFuncMatcher struct {
		s string
		f matcherFunc
	}
)

var _ gomock.Matcher = &matcherFuncMatcher{}

func NewMatcherFunc(s string, f matcherFunc) gomock.Matcher {
	return &matcherFuncMatcher{
		s: s,
		f: f,
	}
}

func (f *matcherFuncMatcher) Matches(x interface{}) bool {
	return f.f(x)
}

func (f *matcherFuncMatcher) String() string {
	return f.s
}
