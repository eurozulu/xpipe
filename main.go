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

	// convert arg parms to datastreams
	dss := make([]xpipe.DataStreams, len(pipeNames))
	for i, pn := range pipeNames {
		ds := xpipe.NewDataStreams(pn)
		if err := ds.Open(); err != nil {
			logger.Error("Failed to open parameter %d, '%s' as a network address.  %v", i, pn, err)
		}
		dss[i] = ds
	}

	// build the pipeline and hook it up to stdout
	var pipeline xpipe.Pipeline
	pipeline.Output = os.Stdout
	pipeline.InitPipeline(len(dss))

	defer func() {
		logger.Debug("closing pipeline")
		if err := pipeline.Close(); err != nil {
			logger.Error("failed to close pipeline %v", err)
			return
		}
	}()

	pipeline.Connect(dss)

	logger.Debug("pipeline opened")

	for {
		select {
		case <-ctx.Done():
			logger.Debug("context completed. %v", ctx.Err())
			return
		case <-sig:
			logger.Debug("shutdown signal received from OS")
			return
		}
	}
}
