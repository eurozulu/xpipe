package xpipe

import (
	"io"
	"os"
)

type ConsoleStreams struct {
}

func (c ConsoleStreams) Open() error {
	return nil
}

func (c ConsoleStreams) Close() error {
	return nil
}

func (c ConsoleStreams) Connect() (io.Reader, io.Writer) {
	return os.Stdin, os.Stdout
}

func (c ConsoleStreams) Disconnect() error {
	return nil
}
