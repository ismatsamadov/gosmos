package main

import (
	"errors"
	"net"
	"sync"
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
