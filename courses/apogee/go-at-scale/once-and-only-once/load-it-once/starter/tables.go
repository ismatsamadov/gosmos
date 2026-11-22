package main

import "sync"

// Tables holds what passplan 2 and the Board read from Port Ellery's file server, each loaded once for every caller.
type Tables struct {
	once sync.Once
	net  *Network
	err  error
}

// Network returns the network, the twenty stations and eighty satellites, loading it the first time any caller asks.
func (t *Tables) Network() (*Network, error) {
	return loadNetwork() // TODO: load once, in t.once.Do, and keep what it returns in t
}
