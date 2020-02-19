package xpipe

import (
	"github.com/eurozulu/xpipe/logger"
	"io"
	"strings"
)

// Streams provides a Connection to a source of stream data.
// specifically, it provides a Reader and Writer to a Piped process, either a network connection or a console process
type Streams interface {

	// Connect to the resource to gets its streams.
	// connect is blocking, until the data stream is available.
	Connect() (Connection, error)

	// Disconnect from the current connection and close the streams
	Disconnect() error

}

// NewStreams creates a new Streams based on the given name.
// if name is a '-', then a ConsoleStreams is created
// If name is in brackets an Inbound Network Streams is created
// Otherwise an outbound network connection is created
func NewStreams(pn string) (Streams, error) {
	if "-" == pn {
		return &ConsoleStreams{}, nil
	}

	if strings.HasPrefix(pn, "(") && strings.HasSuffix(pn, ")") {
		return NewNetAddrListener(pn[1: len(pn) - 1])
	}
	return &NetAddrSender{NetAddr: pn}, nil
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
			logger.Debug("closing Reader stream. %v", err)
			return err
		}
		if l > 0 {
			logger.Debug("writing %d bytes through stream.", l)
			_, err := w.Write(buf[0: l])
			if err != nil {
				logger.Debug("Closing Writer stream. %v", err)
				return err
			}
		}
	}
}