package xpipe

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/xpipe/logger"
	"io"
	"os"
	"sync"
)

// Pipeline is the collection of pipes which make up the complete stream
type Pipeline struct {
	Output io.WriteCloser
	pipes  []*Pipe
	data   []DataStreams
}

func (pl Pipeline) Len() int { return len(pl.pipes) }

func (pl Pipeline) FirstPipe() *Pipe {
	if pl.Len() < 1 {
		return nil
	}
	return pl.pipes[0]
}

func (pl Pipeline) LastPipe() *Pipe {
	if pl.Len() < 1 {
		return nil
	}
	return pl.pipes[pl.Len()-1]
}

// Connect connects the pipeline to the datasources.
func (pl *Pipeline) Connect(streams []DataStreams) {
	var wg sync.WaitGroup
	pl.data = streams

	for i, ds := range pl.data {
		wg.Add(1)
		// Start the connection, ds.Connect will block on netaddr's, until a connection is established.
		go func(p *Pipe, streams DataStreams) {
			wg.Done()
			r, w := streams.Connect()
			p.OnConnect(r, w)
		}(pl.pipes[i], ds)
	}
	wg.Wait() // wait for it all to begin
}

func (pl Pipeline) Close() error {
	var errs bytes.Buffer
	for _, p := range pl.pipes {
		err := p.Close()
		if err != nil {
			errs.WriteString(err.Error())
			errs.WriteString("  ")
		}
	}
	if errs.Len() > 0 {
		return fmt.Errorf("%s", errs.String())
	}
	return nil
}


// InitPipeline initalises the pipeline, creating a new Pipe segment for each of the given datastreams.
func (pl *Pipeline) InitPipeline(size int) {
	pl.pipes = make([]*Pipe, size)

	for i := 0; i < size; i ++ {
		pl.pipes[i] = &Pipe{}
		if i == 0 {
			continue
		}
		in, out := io.Pipe()
		pl.pipes[i-1].Output = out
		pl.pipes[i].Input = in
	}

	if pl.Output == nil {
		logger.Info("pipeline output not set, defaulting to stdout")
		pl.Output = os.Stdout
	}

	// Set the last pipe to pipeline output
	pl.LastPipe().Output = pl.Output
}
