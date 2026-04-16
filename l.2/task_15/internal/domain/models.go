package domain

import (
	"io"
	"os"
)

type Command struct {
	Text         string
	StdIn        string
	StdOut       string
	StdErr       string
	OpenFlagIn   int
	OpenFlagOut  int
	OpenFlagErr  int
	PipeReader   io.ReadCloser
	PipeWriter   io.WriteCloser
	PipeChanKill chan struct{}
	Pipe         bool
	And          bool
	Or           bool
}

type Conditionals struct {
	Subsequence []*Command
	Errs        []error
}

func NewCommand(text string) *Command {
	return &Command{
		Text:        text,
		StdIn:       "",
		StdOut:      "",
		StdErr:      "",
		OpenFlagIn:  os.O_RDONLY,
		OpenFlagOut: os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
		OpenFlagErr: os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
		PipeReader:  nil,
		PipeWriter:  nil,
		Pipe:        false,
		And:         false,
		Or:          false,
	}
}

func NewConditionals() *Conditionals {
	return &Conditionals{
		Subsequence: make([]*Command, 0),
		Errs:        make([]error, 0),
	}
}
