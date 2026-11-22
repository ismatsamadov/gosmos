package main

import (
	"bufio"
	"context"
)

// pass is one contact on the pass line, from its AOS to its LOS.
type pass struct {
	r *run
	c Contact

	answered bool // the station answered the call with its banner
}

// call calls the pass's station, reads its banner and holds the call until ctx ends, then
// hangs up; a failed call or a wrong banner hangs up at once.
func (p *pass) call(ctx context.Context) {
	conn, err := p.r.w.Dial(ctx, p.c.Station)
	if err != nil {
		return
	}
	defer conn.Close()
	br := bufio.NewReader(conn)
	if err := readBanner(br, p.c.Station); err != nil {
		return
	}
	p.answered = true
	<-ctx.Done()
}

// outcome says how the pass ended: flown once its station answered, failed no_answer
// otherwise.
func (p *pass) outcome() Outcome {
	if !p.answered {
		return Outcome{"failed", "no_answer"}
	}
	return Outcome{Kind: "flown"}
}
