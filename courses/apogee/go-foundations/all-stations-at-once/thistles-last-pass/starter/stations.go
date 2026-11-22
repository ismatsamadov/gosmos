package main

// Heard is one beacon heard at one station.
type Heard struct {
	Seq     int
	Station string // station code
}

// listen plays back what one station heard of Thistle's last pass: each beacon, then the end.
func listen(station string, seqs ...int) <-chan Heard {
	feed := make(chan Heard)
	go func() {
		for _, seq := range seqs {
			feed <- Heard{Seq: seq, Station: station}
		}
		close(feed)
	}()
	return feed
}
