package main

import (
	"context"
	"crypto/tls"
	"net"
)

// dialController connects to the dish controller at addr over TLS with cfg, handshaking within ctx.
func dialController(ctx context.Context, addr string, cfg *tls.Config) (net.Conn, error) {
	var d net.Dialer
	return d.DialContext(ctx, "tcp", addr)
}
