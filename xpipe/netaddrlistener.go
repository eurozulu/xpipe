package xpipe

import (
	"github.com/eurozulu/xpipe/logger"
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
	conn       net.Conn
}

func (n *NetAddrListener) Connect() (Connection, error) {
	if n.conn != nil {
		return n.conn, nil
	}
	logger.Info("Listening for inbound connection on %v", n.NetAddr)
	c, err := n.listenPort.Accept()
	if err != nil {
		// Check if its connection closing
		if strings.HasSuffix(err.Error(), "use of closed network connection") {
			logger.Debug("network connection on %s closing", n.NetAddr)
			return nil, nil
		}
		return nil, err
	}
	n.conn = c

	logger.Info("inbound network connection on '%s' from '%s'", n.NetAddr, n.conn.RemoteAddr().String())
	return n.conn, nil
}

// Disconnect any active connection
func (n *NetAddrListener) Disconnect() error {
	if n.conn != nil {
		c := n.conn
		n.conn = nil
		return c.Close()
	}
	return nil
}

func NewNetAddrListener(netAddr string) (*NetAddrListener, error) {
	logger.Debug("opening inbound network port on %s", netAddr)
	l, err := net.Listen(Network, netAddr)
	if err != nil {
		return nil, err
	}
	return &NetAddrListener{
		NetAddr:    l.Addr().String(),
		listenPort: l,
	}, nil
}
