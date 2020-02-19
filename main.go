package main

import (
	"context"
	"fmt"
	"github.com/eurozulu/xpipe/logger"
	"github.com/eurozulu/xpipe/xpipe"
	"os"
	"os/signal"
	"path"
	"sync"
	"time"
)

func main() {
	args := Arguments{Args: os.Args[1:]}
	pipeNames := args.Parameters()
	if len(pipeNames) < 1 {
		fmt.Printf("%s [OPTIONS]... SOURCE [TARGET]...\n", path.Base(os.Args[0]))
		os.Exit(-1)
	}

	// Set up logging level
	lvl := logger.WARN
	vb, ok := args.Flag("verbose", "v")
	if ok {
		lvl = logger.INFO
		if vb != "" {
			lvl = logger.ParseLogLevel(vb)
		}
	}
	logger.DefaultLog.Level = lvl
	logger.Debug("Loglevel: %d", lvl)

	nw, ok := args.Flag("network")
	if ok {
		logger.Debug("network type set to %s", nw)
		xpipe.Network = nw
	}

	to, ok := args.FlagDuration("timeout")
	if ok {
		logger.Debug("timeout set to %s", to.String())
		xpipe.NetTimeout = time.Duration(to)
	}

	var ctx context.Context
	var cnl context.CancelFunc
	if xpipe.NetTimeout == 0 {
		logger.Warn("network timeout set to zero")
		ctx, cnl = context.WithCancel(context.Background())
	} else {
		ctx, cnl = context.WithTimeout(context.Background(), xpipe.NetTimeout)
	}
	defer cnl()

	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, os.Kill, os.Interrupt)

	pipeNames = append(pipeNames, "-") // out final stdOut as last pipe

	strs, err := buildStreams(pipeNames)
	if err != nil {
		logger.Error("Error: Invalid arguments. %v", err)
		return
	}

	pipeline := make(xpipe.Pipeline, len(strs))
	chConns := openPipeline(strs)

	for {
		select {
		case <-ctx.Done():
			logger.Debug("context completed. %v", ctx.Err())
			return
		case <-sig:
			logger.Debug("shutdown signal received from OS")
			return
		case ci, ok := <-chConns:
			if !ok {
				logger.Debug("pipeline connections are completed. All pipes are open.")
			}
			if ci.Err != nil {
				logger.Error("Error: %v", ci.Err)
				return
			}
			if ci.Index >= len(pipeline) {
				logger.Error("unexpected connection index of %d. Out of range of the pipeline length of %d", ci.Index, len(pipeline))
				return
			}
			if pipeline[ci.Index] != nil {
				logger.Error("unexpected second connection at index %d.", ci.Index)
				return
			}

			pipeline[ci.Index] = ci.Connected
			go connectToPeers(pipeline.Prev(ci.Index), pipeline[ci.Index], pipeline.Next(ci.Index))
		}
	}
}

type connectedIndex struct {
	Connected xpipe.Connection
	Err       error
	Index     int
}

func connectToPeers(prev xpipe.Connection, nu xpipe.Connection, next xpipe.Connection) {
	defer func(c xpipe.Connection) {
		if c != nil {
			if err := c.Close(); err != nil {
				logger.Error("Error: failed to close connection %v", err)
			}
		}
	}(nu)

	// hook upto the next connection in the pipeline, if present
	if next != nil {
		if err := xpipe.
			CopyStream(nu, next); err != nil {
			logger.Error("stream connection closed %v.", err)
			return
		}
	}
	// hook upto the prev connection in the pipeline, if present
	if prev != nil {
		if err := xpipe.CopyStream(prev, nu); err != nil {
			logger.Error("stream connection closed %v.", err)
			return
		}
	}
}

func openPipeline(strs []xpipe.Streams) <-chan *connectedIndex {
	conns := make(chan *connectedIndex, len(strs))
	var wg sync.WaitGroup

	for i, str := range strs {
		wg.Add(1)
		go func(index int, st xpipe.Streams, ch chan<- *connectedIndex) {
			defer wg.Done()

			ci := &connectedIndex{
				Index: index,
			}
			cnt, err := st.Connect()
			if err != nil {
				ci.Err = err
			} else {
				ci.Connected = cnt
			}
			ch <- ci
		}(i, str, conns)
	}
	go func(ch chan<- *connectedIndex, w sync.WaitGroup) {
		w.Wait()
		close(ch)
	}(conns, wg)

	logger.Debug("pipeline opened")
	return conns
}

// build the full pipe line, to allow listeners to fire up, prior to streaming.
func buildStreams(addrs []string) ([]xpipe.Streams, error) {
	strs := make([]xpipe.Streams, len(addrs))
	for i, addr := range addrs {
		str, err := xpipe.NewStreams(addr)
		if err != nil {
			return nil, fmt.Errorf("Failed to '%s' as a network address.  %v", addr, err)
		}
		strs[i] = str
	}
	return strs, nil
}
