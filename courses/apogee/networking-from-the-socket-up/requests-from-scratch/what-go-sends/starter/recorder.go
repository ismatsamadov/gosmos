package main

import (
	"context"
	"net"
	"net/http"
	"time"
)

// askForList sends GET http://127.0.0.1/recordings from net/http's own client, as the
// recorder and the collector do, through a connection to the listener at addr. The
// channel it returns closes once the answer has been read.
func askForList(addr string) <-chan struct{} {
	done := make(chan struct{})
	transport := &http.Transport{DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, network, addr)
	}}
	client := &http.Client{Transport: transport, Timeout: time.Second}
	go func() {
		defer close(done)
		defer transport.CloseIdleConnections()
		res, err := client.Get("http://127.0.0.1/recordings")
		if err == nil {
			res.Body.Close()
		}
	}()
	return done
}
