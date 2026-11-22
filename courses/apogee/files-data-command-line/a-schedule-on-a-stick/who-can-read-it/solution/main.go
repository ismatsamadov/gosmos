package main

import "os"

// ownerOnly reports whether no group or other permission bits are set.
func ownerOnly(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().Perm()&0077 == 0, nil
}
