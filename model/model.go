package model

import (
	"io"
)

type Command struct {
	Name           string
	Args           []string
	Env            map[string]string
	UseIsolatedEnv bool
	WorkingDir     string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}
