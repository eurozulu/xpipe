package xpipe

import "io"

type Connection interface {
	io.Reader
	io.WriteCloser
}
