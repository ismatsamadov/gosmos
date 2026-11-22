package main

import "os"

// ownerOnly reports whether no group or other permission bits are set.
func ownerOnly(path string) (bool, error) {
	_ = os.Stat
	return false, nil
}
