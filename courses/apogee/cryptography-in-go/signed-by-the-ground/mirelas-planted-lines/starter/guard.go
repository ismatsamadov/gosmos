package main

import (
	"errors"
	"fmt"
)

var (
	// ErrForged refuses a command whose signature doesn't check out, and never says more.
	ErrForged = errors.New("command failed authentication")
	// ErrReplay refuses a command whose sequence number isn't above the last one accepted.
	ErrReplay = errors.New("replayed sequence number")
)

// Guard refuses forged and replayed commands; its zero value is ready to use.
type Guard struct {
	last map[string]uint64 // satellite ID → last accepted sequence number
}

// Accept checks c's signature under key and its sequence number, then records the number.
func (g *Guard) Accept(c Command, key []byte) error {
	text, err := canonical(c)
	if err != nil {
		return err
	}
	if !verify(text, c.Sig, key) {
		return ErrForged
	}
	if last := g.last[c.SatID]; c.Seq <= last {
		return fmt.Errorf("%w: %d, last %d", ErrReplay, c.Seq, last)
	}
	if g.last == nil {
		g.last = make(map[string]uint64)
	}
	g.last[c.SatID] = c.Seq
	return nil
}
