package main

import (
	"context"
	"fmt"
	"net"
)

// connect gives the link a new connection: it redials, reads the banner and says HELLO,
// and while the handshake fails as a lost line, tries it all again after redialDelay,
// counting failed handshakes in a row as redial counts failed dials. On a WELCOME it
// learns the controller's count and calls resume, so afterwards l.inflight is nil if the
// command in flight has already run. It returns ctx's error when ctx ends first.
func (l *Link) connect(ctx context.Context) error {
	for failed := 1; ; failed++ {
		conn, err := l.redial(ctx)
		if err != nil {
			return err
		}
		w, err := l.greet(conn)
		if err == nil {
			l.learn(w)
			l.resume(w)
			return nil
		}
		conn.Close()
		if !lost(err) {
			return err
		}
		if err := l.Wait(ctx, redialDelay(failed, l.Rand)); err != nil {
			return err
		}
	}
}

// greet reads the banner on conn and says HELLO; on a WELCOME the client becomes the
// link's.
func (l *Link) greet(conn net.Conn) (Welcome, error) {
	c, err := newClient(conn)
	if err != nil {
		return Welcome{}, err
	}
	c.Timeout = l.Timeout
	w, err := c.hello()
	if err != nil {
		return Welcome{}, err
	}
	l.client = c
	return w, nil
}

// settle judges the controller's answer to cmd: an ACK of its seq is nil, and so is a
// NAK duplicate, which means it has already run; any other NAK of its seq is its
// *NakError, and anything else an error.
func settle(cmd, answer Frame) error {
	if answer.Seq == cmd.Seq {
		switch answer.Type {
		case TypeAck:
			return nil
		case TypeNak:
			if nak := nakError(answer); nak.Code != 1 {
				return nak
			}
			return nil // duplicate: it already ran
		}
	}
	return fmt.Errorf("answer to seq %d: type %#02x, seq %d", cmd.Seq, answer.Type, answer.Seq)
}
