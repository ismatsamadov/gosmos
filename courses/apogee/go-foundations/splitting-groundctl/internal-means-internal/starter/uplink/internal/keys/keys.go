// Package keys holds the uplink signing keys. Only uplink may import it.
package keys

// Key is a signing key and the ID ops quote for it.
type Key struct {
	ID     string
	Secret []byte
}

// Current returns the key that signs commands now.
func Current() Key {
	return Key{ID: "k-2027-03", Secret: []byte("test key, never flown")}
}
