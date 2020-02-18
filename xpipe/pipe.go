package xpipe

import (
	"fmt"
	"github.com/eurozulu/xpipe/logger"
	"io"
	"sync"
)

// Pipe is a segment of the pipeline representing two streams, an Input and Output.
// These stream, when set, are copied into a connecting set of streams when those streams become available.
type Pipe struct {
	// Input is the reader this pipe will read from to write out to the external connection
	Input io.ReadCloser

	// Output is the writer this pipe will write the input of an external connection
	Output io.WriteCloser
}

// OnConnect copies the given streams into the respective pipe streams.
// the given Reader is read into the pipe's Output, the pipes Input is read into given Writer
// If with the given stream or any of the Pipe streams are nil, the respective copy will not take place.
func (c Pipe) OnConnect(r io.Reader, w io.Writer) {
	var wg sync.WaitGroup

	// Copy given reader into pipes output
	if r != nil && c.Output != nil {
		wg.Add(1)
		go func() {
			CopyStream(r, c.Output)
			wg.Done()
			logger.Debug("Pipe output stream closing")
		}()
	}
	// Copy pipe input to the given output
	if w != nil && c.Input != nil {
		wg.Add(1)
		go func() {
			CopyStream(c.Input, w)
			wg.Done()
			logger.Debug("Pipe input stream closing")
		}()
	}
	wg.Wait()
}

func (c Pipe) Close() error {
	var errIn error
	var errOut error

	if c.Input != nil {
		errIn = c.Input.Close()
	}
	if c.Output != nil {
		errOut = c.Output.Close()
	}
	if errIn == nil {
		return errOut
	}
	if errOut == nil {
		return errIn
	}
	return fmt.Errorf("%v %v", errIn, errOut)
}

// CopyStream copies the given Readers content, into the given writer until EOF or other error occurs.
func CopyStream(r io.Reader, w io.Writer) error {
	if r == nil || w == nil {
		return nil
	}
	buf := make([]byte, 1024)
	for {
		l, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		if l > 0 {
			_, err := w.Write(buf[0: l])
			if err != nil {
				if err != io.EOF {return err}
				return nil
			}
		}
	}
}
