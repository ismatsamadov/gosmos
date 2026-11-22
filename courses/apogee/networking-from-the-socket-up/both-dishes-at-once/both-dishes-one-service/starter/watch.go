package main

import (
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

// serve accepts connections on ln until it is closed, running handle for each one in a
// goroutine of its own, counted by wg. It returns nil once ln is closed, and any other Accept
// error as it is.
func serve(ln net.Listener, wg *sync.WaitGroup, handle func(net.Conn)) error {
	for {
		c, err := ln.Accept()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		if err != nil {
			return err
		}
		wg.Go(func() { handle(c) })
	}
}

// Broadcast writes line and a newline to every watcher, giving each write w.Timeout, and
// removes each one whose write fails.
func (w *Watchers) Broadcast(line string) {
	var conns []net.Conn
	w.mu.Lock()
	for c := range w.conns {
		conns = append(conns, c)
	}
	w.mu.Unlock()
	for _, c := range conns {
		c.SetWriteDeadline(time.Now().Add(w.Timeout))
		if _, err := io.WriteString(c, line+"\n"); err != nil {
			w.Remove(c)
		}
	}
}
