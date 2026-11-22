package main

import (
	"fmt"
	"strings"
)

// String names the key and hides its secret, for %v, %+v and %s.
func (k Key) String() string {
	return k.ID + " (secret hidden)"
}

// GoString names the key and hides its secret, for %#v.
func (k Key) GoString() string {
	return fmt.Sprintf("Key{ID: %q, Secret: hidden}", k.ID)
}

// String lists the ring's key IDs in activation order, and no secret.
func (r *KeyRing) String() string {
	ids := make([]string, len(r.versions))
	for i, v := range r.versions {
		ids[i] = v.ID
	}
	return "key ring: " + strings.Join(ids, ", ")
}

// GoString describes the ring as String does, for %#v.
func (r *KeyRing) GoString() string {
	return r.String()
}
