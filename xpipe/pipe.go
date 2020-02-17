package xpipe

import (
	"bufio"
	"context"
	"github.com/eurozulu/xpipe/logger"
	"io"
)

// Block size of each chunk of streamed data
const DefaultPacketSize = 1024
// Packet is a chunk of data being sent through the stream
type Packet []byte


// Pipeline is the colelction of pipes which make up the complete stream
type Pipeline []Pipe


// Pipe is a segment of the Pipeline which can provide streams to read and write bytes.
type Pipe interface {
	InputStream(ctx context.Context) (io.Reader, error)
	OutputStream(ctx context.Context) (io.Writer, error)
}

// NewPipe createa a new pipe based on the given pipe name.
// if name is a '-', then a ConsolePipe is created, otherwise a NetPipe is created with the given name used as the netAddress.
func NewPipe(pn string) Pipe {
	if "-" == pn {
		return &ConsolePipe{}
	}
	return &NetPipe{NetAddr:pn}
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

// ReadStream reads the given Reader into the resulting channel of Packets.
// Non blocking, channel remains open as long as the reader is open and until context is Done
func ReadStream(ctx context.Context, in io.Reader) (<-chan Packet, error) {
	ch := make(chan Packet)
	go func(c context.Context, r io.Reader, pc chan<- Packet) {
		defer close(pc)
		br := bufio.NewReader(r)
		buf := make([]byte, DefaultPacketSize)

		for {
			l, err := br.Read(buf)
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
	}(ctx, in, ch)
	return ch, nil
}

// OpenPipeline opens the given pipeline for reading.
// Beginning at the first, each Pipe has its inputstream opened and passed into a new channel of Packets.
// The next pipe, should it exist, also has its input stream opened, and the previous channel is passed into
// its output stream, connecting the putput of the previous Pipe, to the input of the current pipe.
// This continues along the pipeline until the last segment, from which the output channel is returned.
// All resulting functions running wach pipe are controlled by the given context.
func OpenPipeline(ctx context.Context, pl Pipeline) (<-chan Packet, error) {
	logger.Log(logger.DEBUG, "opening pipeline length %d:\n", len(pl))

	var pkts <-chan Packet
	for _, p := range pl {
		in, err := p.InputStream(ctx)
		if err != nil {
			return nil, err
		}

		ch, err := ReadStream(ctx, in)
		if err != nil {
			return nil, err
		}

		// Connect new Pipe to any previous pipe
		if pkts != nil {
			out, err := p.OutputStream(ctx)
			if err != nil {
				return nil, err
			}
			go WriteStream(ctx, out, pkts)
		}
		pkts = ch
	}
	logger.Log(logger.DEBUG, "pipeline open")
	return pkts, nil
}
