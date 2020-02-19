package xpipe

import (
	"os"
)

type ConsoleEcho struct {
}

func (c ConsoleEcho) Connect() (Connection, error) {
	return &ConsoleEchoConnection{}, nil
}

func (c ConsoleEcho) Disconnect() error {
	return nil
}

type ConsoleEchoConnection struct{}

func (c ConsoleEchoConnection) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (c ConsoleEchoConnection) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (c ConsoleEchoConnection) Close() error {
	// do nothing on console streams
	return nil
}
