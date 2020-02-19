package xpipe

import (
	"github.com/eurozulu/xpipe/logger"
	"net"
)

type NetAddrSender struct {
	NetAddr string
	conn    net.Conn
}

func (n *NetAddrSender) Connect() (Connection, error) {
	if n.conn != nil {
		return n.conn, nil
	}
	logger.Info("Opening connection to %v", n.NetAddr)
	c, err := net.Dial(Network, n.NetAddr)
	if err != nil {
		return nil, err
	}
	logger.Debug("connection to %v opened", c.RemoteAddr())
	n.conn = c

	return n.conn, nil
}

func (n *NetAddrSender) Disconnect() error {
	if n.conn != nil {
		c := n.conn
		n.conn = nil
		return c.Close()
	}
	return nil
}
