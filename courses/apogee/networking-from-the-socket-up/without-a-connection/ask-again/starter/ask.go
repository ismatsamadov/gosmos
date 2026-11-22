package main

import (
	"net"
	"time"
)

// ask asks the recorder question n from pc and returns its answer's text,
// asking again, under the same number, each time wait passes without one,
// up to tries times.
func ask(pc net.PacketConn, recorder net.Addr, n uint32, wait time.Duration, tries int) (string, error) {
	if _, err := pc.WriteTo(question(n), recorder); err != nil {
		return "", err
	}
	pc.SetReadDeadline(time.Now().Add(wait))
	// TODO: when the wait passes with no answer, ask again, up to tries times
	return readAnswer(pc, recorder, n)
}
