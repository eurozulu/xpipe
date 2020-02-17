package xpipe

import (
	"context"
	"github.com/eurozulu/xpipe/logger"
	"io"
	"net"
	"strings"
	"time"
)

const defaultNetwork = "tcp"
const defaultNetworkTimeout = time.Second * 30

var Network string
var NetTimeout time.Duration

func init() {
	if Network == "" {
		Network = defaultNetwork
	}
	if NetTimeout == 0 {
		NetTimeout = defaultNetworkTimeout
	}
}

type NetPipe struct {
	NetAddr string

	listenPort net.Listener
}

func (n NetPipe) InputStream(ctx context.Context) (io.Reader, error) {
	l, err := n.listener(ctx)
	if err != nil {
		return nil, err
	}

	logger.Log(logger.INFO, "waiting for inbound network connection on %s", n.NetAddr)
	c, err := l.Accept()
	if err != nil {
		// Check if its connection closing
		if strings.HasSuffix(err.Error(), "use of closed network connection") {
			logger.Log(logger.INFO, "network connection  on %s closing", n.NetAddr)
			return nil, nil
		}
		return nil, err
	}
	logger.Log(logger.DEBUG, "inbound network connection on '%s' from '%v'", n.NetAddr, c.RemoteAddr())

	go func(cx context.Context, con net.Conn) {
		defer func() {
			logger.Log(logger.DEBUG, "closing network connection")
			con.Close()
		}()
		<-cx.Done()
	}(ctx, c)
	return c, nil
}

func (n NetPipe) OutputStream(ctx context.Context) (io.Writer, error) {
	logger.Log(logger.INFO, "Opening connection to %v", n.NetAddr)
	c, err := net.Dial(Network, n.NetAddr)
	if err != nil {
		return nil, err
	}
	logger.Log(logger.DEBUG, "connection to %v opened", c.RemoteAddr())

	go func(cx context.Context, con net.Conn) {
		defer con.Close()
		<-cx.Done()
	}(ctx, c)
	return c, nil
}

func (n *NetPipe) listener(ctx context.Context) (net.Listener, error) {
	if n.listenPort == nil {
		logger.Log(logger.DEBUG, "opening network connection to %s", n.NetAddr)
		l, err := net.Listen(Network, n.NetAddr)
		if err != nil {
			return nil, err
		}
		logger.Log(logger.DEBUG, "network connection to %s opened", n.NetAddr)

		n.listenPort = l
		go func() {
			defer func() {
				logger.Log(logger.DEBUG, "closing network connection to %s", n.NetAddr)
				l.Close()
			}()
			<-ctx.Done()
		}()
	}
	return n.listenPort, nil
}
