package main

// Key is a signing key and the ID ops quote for it.
type Key struct {
	ID     string
	Secret []byte
}
