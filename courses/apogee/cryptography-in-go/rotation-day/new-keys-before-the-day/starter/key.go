package main

import "fmt"

// Key is a signing key and the ID ops quote for it.
type Key struct {
	ID     string
	Secret []byte
}

// String names the key and hides its secret, for %v, %+v and %s.
func (k Key) String() string {
	return k.ID + " (secret hidden)"
}

// GoString names the key and hides its secret, for %#v.
func (k Key) GoString() string {
	return fmt.Sprintf("Key{ID: %q, Secret: hidden}", k.ID)
}
