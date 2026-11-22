package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"
)

// pass is one contact on the pass line, from its AOS to its LOS.
type pass struct {
	r *run
	c Contact

	welcomed bool // a WELCOME has come on one of its calls
	nakked   bool // the station refused its HELLO with a code other than stowed
	skipped  bool // the station refused its first HELLO: its dishes were stowed
	cut      bool // the station cut it: BYE stowed, or a redial refused as stowed
	lost     bool // its line is lost and not back
}

// callEnd says how one call of a pass ended.
type callEnd int

const (
	callOver   callEnd = iota // the pass is over: LOS, BYE, or an answer that ends it
	callLost                  // the line was lost after a WELCOME: dial again at once
	callFailed                // no WELCOME on this call: dial again after a wait
)

// received is a frame the reader read, or the error that ended its reading.
type received struct {
	f   Frame
	err error
}

// fly calls the station and keeps the pass on the line until ctx ends: a lost line is
// dialed again at once, and every failed call waits a little longer before the next.
func (p *pass) fly(ctx context.Context) {
	failures := 0
	for {
		switch p.call(ctx) {
		case callOver:
			return
		case callLost:
			failures = 0
		case callFailed:
			failures++
			if !sleep(ctx, redialDelay(failures)) {
				return
			}
		}
	}
}

// redialDelay is the wait after the n-th failed call in a row: 1 s doubled for each
// earlier failure, at most 30 s, times a factor drawn from [0.5, 1.0).
func redialDelay(n int) time.Duration {
	d := 30 * time.Second
	if n < 6 {
		d = time.Second << (n - 1)
	}
	return time.Duration(float64(d) * (0.5 + rand.Float64()/2))
}

// call makes one call to the pass's station: the banner, HELLO, then every frame until the
// call or the pass ends. When ctx ends first, it hangs up at that moment, which ends any
// read still waiting on the call.
func (p *pass) call(ctx context.Context) callEnd {
	conn, err := p.r.w.Dial(ctx, p.c.Station)
	if err != nil {
		if ctx.Err() != nil {
			return callOver
		}
		return callFailed
	}
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	defer stop()
	defer conn.Close()
	br := bufio.NewReader(idleReader{conn})
	if err := readBanner(br, p.c.Station); err != nil {
		if ctx.Err() != nil {
			return callOver
		}
		return callFailed
	}
	hello := Frame{Type: TypeHello, Body: fmt.Appendf(nil, "%s %d", p.c.ID, p.c.Dish)}
	if err := writeFrame(conn, hello); err != nil {
		if ctx.Err() != nil {
			return callOver
		}
		return callFailed
	}
	in := make(chan received)
	quit, done := make(chan struct{}), make(chan struct{})
	go readFrames(br, in, quit, done)
	defer func() {
		conn.Close()
		close(quit)
		<-done
	}()
	return p.talk(ctx, in)
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

// talk is one call after its HELLO: it takes each frame until the call or the pass ends,
// and says how the call ended.
func (p *pass) talk(ctx context.Context, in <-chan received) callEnd {
	welcomed := false // on this call
	for {
		select {
		case <-ctx.Done():
			return callOver
		case got := <-in:
			if errors.Is(got.err, errBadCRC) {
				continue
			}
			if got.err != nil {
				if ctx.Err() != nil {
					return callOver
				}
				if !welcomed {
					return callFailed
				}
				p.lost = true
				return callLost
			}
			switch f := got.f; f.Type {
			case TypeWelcome:
				welcomed, p.welcomed, p.lost = true, true, false
			case TypeNak:
				p.refused(f)
				return callOver
			case TypeTLM:
				p.telemetry(f.Body)
			case TypeBye:
				p.cut = string(f.Body) == "stowed"
				return callOver
			}
		}
	}
}

// refused takes the station's NAK to a HELLO: stowed skips a pass never welcomed and cuts
// one that was; any other code refuses it.
func (p *pass) refused(f Frame) {
	switch {
	case len(f.Body) == 0 || f.Body[0] != NakStowed:
		p.nakked = true
	case p.welcomed:
		p.cut = true
	default:
		p.skipped = true
	}
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
