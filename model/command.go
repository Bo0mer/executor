package executor

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
