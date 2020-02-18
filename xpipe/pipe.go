package xpipe

import (
	"context"
	"fmt"
	"github.com/eurozulu/xpipe/logger"
	"io"
	"strings"
	"sync"
)

// Block size of each chunk of streamed data
const DefaultPacketSize = 1024

// Packet is a chunk of data being sent through the stream
type Packet []byte

// Pipeline is the colelction of pipes which make up the complete stream
type Pipeline []Pipe

// Pipe is a segment of the Pipeline which can provide streams to read and write bytes.
type Pipe interface {
	InputStream() (io.Reader, error)
	OutputStream() (io.Writer, error)
	Close() error
}

// NewPipe createa a new pipe based on the given pipe name.
// if name is a '-', then a ConsolePipe is created, otherwise a NetPipe is created with the given name used as the netAddress.
func NewPipe(pn string) Pipe {
	if "-" == pn {
		return &ConsolePipe{}
	}
	if strings.HasPrefix(pn, "=>") {
		return &NetPipeIn{NetAddr: pn[2:]}
	}
	return &NetPipeOut{NetAddr: pn}
}

// WriteStream writes the given channel of packets, out to the given Writer.
// Blocks until the packets channel is closed or context is Done
func WriteStream(ctx context.Context, out io.Writer, pkts <-chan Packet) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case pkt, ok := <-pkts:
			if !ok {
				return nil
			}
			if _, err := out.Write(pkt); err != nil {
				return err
			}
		}
	}
}

func OpenPipe(ctx context.Context, pipe Pipe) <-chan Packet {
	ch := make(chan Packet)
	go func(c context.Context, p Pipe, pc chan<- Packet) {
		defer close(pc)
		in, err := p.InputStream()
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			return
		}
		buf := make([]byte, DefaultPacketSize)
		for {
			l, err := in.Read(buf)
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				return
			}
			if l == 0 {
				// Check if ctx has been canceled
				if c.Err() != nil {
					return
				}
				continue
			}
			select {
			case <-c.Done():
				return
			case pc <- buf[0:l]:
			}
		}
	}(ctx, pipe, ch)
	return ch
}

// OpenPipeline opens the given pipeline for reading.
// results in the output of the last pipe in the pipeline.  All prev pipes are being read and the result written to the following pipe.
func OpenPipeline(ctx context.Context, pl Pipeline) (<-chan Packet, error) {
	logger.Log(logger.DEBUG, "opening pipeline length %d:\n", len(pl))
	var pkts <-chan Packet
	for _, p := range pl {
		ch := OpenPipe(ctx, p)

		var wg sync.WaitGroup
		// Connect new Pipe to any previous pipe
		if pkts != nil {
			wg.Add(1)// wait until func starts to get correct assignment of pkts
			go func(cx context.Context, ch <-chan Packet) {
				wg.Done()
				out, err := p.OutputStream()
				if err != nil {
					fmt.Println(err)
					return
				}
				WriteStream(cx, out, ch)
			}(ctx, pkts)
		}
		wg.Wait()
		pkts = ch

	}
	logger.Log(logger.DEBUG, "pipeline open")
	return pkts, nil
}
