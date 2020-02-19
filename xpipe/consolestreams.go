package xpipe

import (
	"os"
)

type ConsoleStreams struct {
}

func (c ConsoleStreams) Connect() (Connection, error) {
	return &ConsoleConnection{}, nil
}

func (c ConsoleStreams) Disconnect() error {
	return nil
}

type ConsoleConnection struct{}

func (c ConsoleConnection) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (c ConsoleConnection) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (c ConsoleConnection) Close() error {
	// do nothing on console streams
	return nil
}
