package xpipe

import (
	"fmt"
	"io"
	"os"
)

// SplitStream is a simple straight pass through, which also writes output to a second stream
type SplitStream struct {
	connection SplitConnections
}

func (c SplitStream) Connect() (Connection, error) {
	if c.connection.in == nil {
		i, o := io.Pipe()
		out := &SplitWriter{
			out:  o,
			out2: os.Stdout,
		}
		c.connection.in = i
		c.connection.out = out
	}
	return &c.connection, nil
}

func (c *SplitStream) Disconnect() error {
	if err := c.connection.Close(); err != nil {
		return err
	}
	c.connection.in = nil
	c.connection.out = nil
	return nil
}

type SplitConnections struct {
	in  *io.PipeReader
	out io.WriteCloser
}

func (cc SplitConnections) Read(p []byte) (n int, err error) {
	return cc.in.Read(p)
}

func (cc SplitConnections) Write(p []byte) (n int, err error) {
	return cc.out.Write(p)
}

func (cc SplitConnections) Close() error {
	var e error
	if nil != cc.in {
		if err := cc.in.Close(); err != nil {
			e = err
		}
		if err := cc.out.Close(); err != nil {
			if e != nil {
				e = fmt.Errorf("%v  %v", e, err)
			} else {
				e = err
			}
		}
	}
	return e
}

type SplitWriter struct {
	out  io.WriteCloser
	out2 io.Writer
}

func (s SplitWriter) Write(p []byte) (n int, err error) {

	if s.out == nil {
		return 0, fmt.Errorf("output has not been set correctly for split streamer")
	}

	l, err := s.out.Write(p)
	if err != nil {
		return l, err
	}
	_, err = s.out2.Write(p)
	if err != nil {
		return 0, err
	}
	return l, nil
}

func (s SplitWriter) Close() error {
	return s.out.Close()
}
