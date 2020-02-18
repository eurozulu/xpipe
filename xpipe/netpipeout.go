package xpipe

import (
	"github.com/eurozulu/xpipe/logger"
	"io"
	"net"
)

type NetPipeOut struct {
	NetAddr string
	conn net.Conn
}

func (n NetPipeOut) Close() error {
	if n.conn != nil {
		return n.conn.Close()
	}
	return nil
}

func (n NetPipeOut) InputStream() (io.Reader, error) {
	return n.connection()
}

func (n NetPipeOut) OutputStream() (io.Writer, error) {
	return n.connection()
}

func (n *NetPipeOut) connection() (net.Conn, error) {
	if n.conn == nil {
		logger.Log(logger.INFO, "Opening connection to %v", n.NetAddr)
		c, err := net.Dial(Network, n.NetAddr)
		if err != nil {
			return nil, err
		}
		logger.Log(logger.DEBUG, "connection to %v opened", c.RemoteAddr())
		n.conn = c
	}
	return n.conn, nil
}
