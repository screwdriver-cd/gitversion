//go:build wireinject
// +build wireinject

package bumper

import (
	"github.com/google/wire"
	"github.com/screwdriver-cd/gitversion/git"
)

var DefaultSet = wire.NewSet(
	wire.Struct(new(DefaultBumper), "*"),
	wire.Bind(new(Bumper), new(*DefaultBumper)),
)

var buildSet = wire.NewSet(
	DefaultSet,
	git.DefaultSet,
)

func NewBumper() Bumper {
	panic(wire.Build(buildSet))
}
