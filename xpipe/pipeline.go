package xpipe

// Pipeline is the collection of connections
type Pipeline []Connection

func (pl Pipeline) First() Connection {
	if len(pl) < 1 {
		return nil
	}
	return pl[0]
}

func (pl Pipeline) Last() Connection {
	last := len(pl) - 1
	if last < 0 {
		return nil
	}
	return pl[last]
}

func (pl Pipeline) Next(index int) Connection {
	next := index + 1
	if next < 0 || next >= len(pl) {
		return nil
	}
	return pl[next]
}

func (pl Pipeline) Prev(index int) Connection {
	prev := index - 1
	if prev < 0 || prev >= len(pl) {
		return nil
	}
	return pl[prev]
}