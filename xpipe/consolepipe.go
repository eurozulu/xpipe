package xpipe

import (
	"io"
	"os"
)

type ConsolePipe struct {
}

func (c ConsolePipe) Close() error {
	// go ahead and do nothing
	return nil
}

func (c ConsolePipe) InputStream() (io.Reader, error) {
	return os.Stdin, nil
}

func (c ConsolePipe) OutputStream() (io.Writer, error) {
	return os.Stdout, nil
}
