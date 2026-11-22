package main

import (
	"errors"
	"net"
	"os"
	"time"
)

// ask asks the recorder question n from pc and returns its answer's text,
// asking again, under the same number, each time wait passes without one,
// up to tries times.
func ask(pc net.PacketConn, recorder net.Addr, n uint32, wait time.Duration, tries int) (string, error) {
	var err error
	for range tries {
		if _, err := pc.WriteTo(question(n), recorder); err != nil {
			return "", err
		}
		pc.SetReadDeadline(time.Now().Add(wait))
		var text string
		text, err = readAnswer(pc, recorder, n)
		if !errors.Is(err, os.ErrDeadlineExceeded) {
			return text, err
		}
	}
	return "", err
}
