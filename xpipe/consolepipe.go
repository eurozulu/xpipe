package xpipe

import (
	"context"
	"io"
	"os"
)

type ConsolePipe struct {
}

func (c ConsolePipe) InputStream(ctx context.Context) (io.Reader, error) {
	return os.Stdin, nil
}

func (c ConsolePipe) OutputStream(ctx context.Context) (io.Writer, error) {
	return os.Stdout, nil
}
