package xpipe

import (
	"bytes"
	"github.com/eurozulu/xpipe/logger"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

const defaultNetwork = "tcp"
const defaultNetworkTimeout = time.Second * 30

var Network string
var NetTimeout time.Duration
var lock sync.Mutex

func init() {
	if Network == "" {
		Network = defaultNetwork
	}
	if NetTimeout == 0 {
		NetTimeout = defaultNetworkTimeout
	}
}

// NetPipeIn opens a local network port and waits for inbound connections.
// On receiving a new connection, the inbound stream is piped to the output and
// the input stream for this pipe is piped back to the new connection.
type NetPipeIn struct {
	NetAddr string

	listenPort net.Listener
	out        PlaceHolderWriter
}

func (n *NetPipeIn) Close() error {
	if n.listenPort != nil {
		l := n.listenPort
		n.listenPort = nil
		return l.Close()
	}
	return nil
}

func (n NetPipeIn) InputStream() (io.Reader, error) {
	return n.connection()
}

func (n NetPipeIn) OutputStream() (io.Writer, error) {
	return n.out, nil
}

func (n *NetPipeIn) connection() (net.Conn, error) {
	lock.Lock()
	defer lock.Unlock()

	l, err := n.listener()
	if err != nil {
		return nil, err
	}

	logger.Log(logger.INFO, "Listening for inbound connection on %v", n.NetAddr)
	c, err := l.Accept()
	if err != nil {
		// Check if its connection closing
		if strings.HasSuffix(err.Error(), "use of closed network connection") {
			logger.Log(logger.INFO, "network connection on %s closing", n.NetAddr)
			return nil, nil
		}
		return nil, err
	}
	logger.Log(logger.DEBUG, "inbound network connection on '%s' from '%v'", n.NetAddr, c.RemoteAddr())

	_, err = n.out.SetOut(c)
	return c, nil
}

func (n *NetPipeIn) listener() (net.Listener, error) {
	if n.listenPort == nil {
		logger.Log(logger.DEBUG, "opening network connection to %s", n.NetAddr)
		l, err := net.Listen(Network, n.NetAddr)
		if err != nil {
			return nil, err
		}
		logger.Log(logger.DEBUG, "network connection to %s opened", n.NetAddr)

		n.listenPort = l
	}
	return n.listenPort, nil
}

type PlaceHolderWriter struct {
	Out    io.Writer
	buffer bytes.Buffer
}

func (w *PlaceHolderWriter) SetOut(out io.Writer) (int, error) {
	var n int
	var err error
	if w.buffer.Len() > 0 {
		n, err = out.Write(w.buffer.Bytes())
		w.buffer.Reset()
	}
	w.Out = out
	return n, err
}

func (w PlaceHolderWriter) Write(p []byte) (n int, err error) {
	if w.Out == nil {
		return w.buffer.Write(p)
	}
	return w.Out.Write(p)
}
