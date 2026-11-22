package main

import "net"

// listenLocal opens a listener for dishd's watch port that only this Mac can reach.
func listenLocal() (net.Listener, error) {
	return net.Listen("tcp", "127.0.0.1:0")
}
