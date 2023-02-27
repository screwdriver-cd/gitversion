//go:build wireinject
// +build wireinject

package git

import (
	"github.com/google/wire"
)

var DefaultSet = wire.NewSet(
	wire.Struct(new(DefaultCmdRunner), "*"),
	wire.Struct(new(DefaultGit), "*"),
	wire.Bind(new(CmdRunner), new(*DefaultCmdRunner)),
	wire.Bind(new(Git), new(*DefaultGit)),
)

var buildSet = DefaultSet

func NewCmdRunner() CmdRunner {
	panic(wire.Build(buildSet))
}

func NewGit() Git {
	panic(wire.Build(buildSet))
}
