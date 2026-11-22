package main

// Key is a signing key and the ID ops quote for it.
type Key struct {
	ID     string
	Secret []byte
}

// String names the key and hides its secret, for %v, %+v and %s.
func (k Key) String() string {
	return k.ID + " (secret hidden)"
}
