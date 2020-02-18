package xpipe

import (
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

// NetAddrListener is a DataStream to an inbound network connection.
// Listening on the NetAddr, as new connection arrive, their stream is read and passed out of the pipe and
// the written to from the input of the pipe.
type NetAddrListener struct {
	// NetAddr is the local address to listen on.  Usually in the form of <ip address>:<port> or simply :<port>
	NetAddr string

	listenPort net.Listener
	conn net.Conn
}

// InitPipeline the local port for listening for inbound connections.
func (n *NetAddrListener) Open() error {
	if n.listenPort == nil {
		logger.Debug("opening inbound network port on %s", n.NetAddr)
		l, err := net.Listen(Network, n.NetAddr)
		if err != nil {
			return err
		}
		logger.Debug("inbound network port on %s opened", n.NetAddr)
		n.listenPort = l
	}
	return nil
}

// Blocks until an inbound connection arrives or has arrived.  Returns the stream from that connection.
func (n *NetAddrListener) Connect() (io.Reader, io.Writer) {
	if n.conn == nil {
		logger.Info("Listening for inbound connection on %v", n.NetAddr)
		c, err := n.listenPort.Accept()
		if err != nil {
			// Check if its connection closing
			if strings.HasSuffix(err.Error(), "use of closed network connection") {
				logger.Info("network connection on %s closing", n.NetAddr)
				return nil, nil
			}
			panic(err)
		}
		n.conn = c
	}
	logger.Debug("inbound network connection on '%s' from '%s'", n.NetAddr, n.conn.RemoteAddr().String())
	return n.conn, n.conn
}

// Disconnect any active connection
func (n *NetAddrListener) Disconnect() error {
	if n.conn != nil {
		c := n.conn
		n.conn = nil
		return c.Close();
	}
	return nil
}

// Close any avtive connection and stop listening on the local port.
func (n *NetAddrListener) Close() error {
	if n.listenPort != nil {
		l := n.listenPort
		n.listenPort = nil
		return l.Close()
	}
	return nil
}
