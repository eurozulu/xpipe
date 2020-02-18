package xpipe

import (
	"io"
	"strings"
)

// DataStreams provides two streams to and from a generic data source.
// specifically, it provides a Reader and Writer to a Piped process, either a network connection or a console process
type DataStreams interface {
	// Opens and prepares it for connections.
	// open is non blocking
	Open() error

	// Connect to the resource to gets its streams.
	// connect is blocking, until the data stream is available.
	Connect() (io.Reader, io.Writer)

	// Disconnect from the current connection and close the streams
	Disconnect() error

	// Close the connection and shutdown any resources its using.
	Close() error
}

// NewDataStreams creates a new DataStreams based on the given name.
// if name is a '-', then a ConsoleStreams is created
// If name begins with '=>' an Inbound Network DataStreams is created
// Otherwise an outbound network connection is created
func NewDataStreams(pn string) DataStreams {
	if "-" == pn {
		return &ConsoleStreams{}
	}
	if strings.HasPrefix(pn, "=>") {
		return &NetAddrListener{NetAddr: pn[2:]}
	}
	return &NetAddrSender{NetAddr: pn}
}

