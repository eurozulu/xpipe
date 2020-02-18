package main

import (
	"context"
	"fmt"
	"github.com/eurozulu/xpipe/logger"
	"github.com/eurozulu/xpipe/xpipe"
	"os"
	"os/signal"
	"path"
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
	logger.Log(logger.DEBUG, "Loglevel: %d", lvl)

	nw, ok := args.Flag("network")
	if ok {
		logger.Log(logger.DEBUG, "network type set to %s", nw)
		xpipe.Network = nw
	}

	to, ok := args.FlagDuration("timeout")
	if ok {
		logger.Log(logger.DEBUG, "timeout set to %s", to.String())
		xpipe.NetTimeout = time.Duration(to)
	}

	var ctx context.Context
	var cnl context.CancelFunc
	if xpipe.NetTimeout == 0 {
		logger.Log(logger.WARN, "network timeout set to zero")
		ctx, cnl = context.WithCancel(context.Background())
	} else {
		ctx, cnl = context.WithTimeout(context.Background(), xpipe.NetTimeout)
	}
	defer cnl()

	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, os.Kill, os.Interrupt)

	// build pipeline
	var pl xpipe.Pipeline
	for _, pn := range pipeNames {
		pl = append(pl, xpipe.NewPipe(pn))
	}

	chPipe, err := xpipe.OpenPipeline(ctx, pl)
	if err != nil {
		logger.Error("failed to open pileline %v", err)
		return
	}

	target := xpipe.NewPipe("-") // final target is standard out
	out, _ := target.OutputStream()

	for {
		select {
		case <-ctx.Done():
			logger.Log(logger.DEBUG, "context completed. %v", ctx.Err())
			return
		case <-sig:
			logger.Log(logger.DEBUG, "shutdown signal received from OS")
			return
		case pkt, ok := <-chPipe:
			if !ok {
				return
			}
			if _, err := out.Write(pkt); err != nil {
				logger.Error("failed to write to output %v", err)
				return
			}
		}
	}
}
