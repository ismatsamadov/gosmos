package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
)

// pass is one contact on the pass line, from its AOS to its LOS.
type pass struct {
	r *run
	c Contact

	welcomed bool // the station took the pass: a WELCOME came
	nakked   bool // the station refused its HELLO with a code other than stowed
	skipped  bool // the station refused its HELLO: its dishes were stowed
	cut      bool // the station cut it: BYE stowed
	lost     bool // its line was lost before LOS
}

// received is a frame the reader read, or the error that ended its reading.
type received struct {
	f   Frame
	err error
}

// call calls the pass's station, reads its banner, says HELLO and takes every frame until
// ctx ends or the station refuses the pass or hangs up; a failed call or a wrong banner
// hangs up at once. When ctx ends first, it hangs up at that moment, which ends any read
// still waiting on the call.
func (p *pass) call(ctx context.Context) {
	conn, err := p.r.w.Dial(ctx, p.c.Station)
	if err != nil {
		return
	}
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	defer stop()
	defer conn.Close()
	br := bufio.NewReader(conn)
	if err := readBanner(br, p.c.Station); err != nil {
		return
	}
	hello := Frame{Type: TypeHello, Body: fmt.Appendf(nil, "%s %d", p.c.ID, p.c.Dish)}
	if err := writeFrame(conn, hello); err != nil {
		return
	}
	in := make(chan received)
	quit, done := make(chan struct{}), make(chan struct{})
	go readFrames(br, in, quit, done)
	defer func() {
		conn.Close()
		close(quit)
		<-done
	}()
	p.talk(ctx, in)
}

// readFrames reads frames until reading fails, sending each on in, damaged ones with
// errBadCRC; it closes done when it returns, and returns at once when quit closes.
func readFrames(br *bufio.Reader, in chan<- received, quit <-chan struct{}, done chan<- struct{}) {
	defer close(done)
	for {
		f, err := readFrame(br)
		select {
		case in <- received{f, err}:
		case <-quit:
			return
		}
		if err != nil && !errors.Is(err, errBadCRC) {
			return
		}
	}
}

// talk is one call after its HELLO: it takes each frame until ctx ends, the station
// refuses the pass or hangs up, or the line is lost.
func (p *pass) talk(ctx context.Context, in <-chan received) {
	for {
		select {
		case <-ctx.Done():
			return
		case got := <-in:
			if errors.Is(got.err, errBadCRC) {
				continue
			}
			if got.err != nil {
				if ctx.Err() == nil {
					p.lost = true // out of step, or gone before LOS without a BYE
				}
				return
			}
			switch f := got.f; f.Type {
			case TypeWelcome:
				p.welcomed = true
			case TypeNak:
				p.refused(f)
				return
			case TypeTLM:
				p.telemetry(f.Body)
			case TypeBye:
				p.cut = string(f.Body) == "stowed"
				return
			}
		}
	}
}

// refused takes the station's NAK to the pass's HELLO: stowed skips the pass, and any
// other code refuses it.
func (p *pass) refused(f Frame) {
	if len(f.Body) > 0 && f.Body[0] == NakStowed {
		p.skipped = true
		return
	}
	p.nakked = true
}

// telemetry takes one TLM frame: a line the strict parser accepts updates the satellite's
// health, and one it refuses is dropped.
func (p *pass) telemetry(body []byte) {
	f, err := parseLine(string(body))
	if err != nil {
		return
	}
	p.r.health.keep(f)
}

// outcome says how the pass ended: the first of these that fits, or flown.
func (p *pass) outcome() Outcome {
	switch {
	case p.skipped:
		return Outcome{"skipped", "stowed"}
	case p.nakked:
		return Outcome{"failed", "refused"}
	case !p.welcomed:
		return Outcome{"failed", "no_answer"}
	case p.cut:
		return Outcome{"failed", "stowed"}
	case p.lost:
		return Outcome{"failed", "line_lost"}
	}
	return Outcome{Kind: "flown"}
}
