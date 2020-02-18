package xpipe

import (
	"github.com/eurozulu/xpipe/logger"
	"io"
	"net"
)

type NetAddrSender struct {
	NetAddr string
	conn net.Conn
}

func (n *NetAddrSender) Connect() (io.Reader, io.Writer) {
	if n.conn == nil {
		logger.Info("Opening connection to %v", n.NetAddr)
		c, err := net.Dial(Network, n.NetAddr)
		if err != nil {
			panic(err)
		}
		logger.Debug("connection to %v opened", c.RemoteAddr())
		n.conn = c
	}
	return n.conn, n.conn
}

func (n *NetAddrSender) Disconnect() error {
	if n.conn != nil {
		c := n.conn
		n.conn = nil
		return c.Close()
	}
	return nil
}

func (n NetAddrSender) Open() error {
	// do nothing.
	return nil
}

func (n NetAddrSender) Close() error {
	return n.Close()
}
