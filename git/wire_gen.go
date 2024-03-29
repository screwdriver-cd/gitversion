// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//go:build !wireinject
// +build !wireinject

package git

import (
	"github.com/google/wire"
)

// Injectors from wire.go:

func NewCmdRunner() CmdRunner {
	defaultCmdRunner := &DefaultCmdRunner{}
	return defaultCmdRunner
}

func NewGit() Git {
	defaultCmdRunner := &DefaultCmdRunner{}
	defaultGit := &DefaultGit{
		CmdRunner: defaultCmdRunner,
	}
	return defaultGit
}

// wire.go:

var DefaultSet = wire.NewSet(wire.Struct(new(DefaultCmdRunner), "*"), wire.Struct(new(DefaultGit), "*"), wire.Bind(new(CmdRunner), new(*DefaultCmdRunner)), wire.Bind(new(Git), new(*DefaultGit)))

var buildSet = DefaultSet
