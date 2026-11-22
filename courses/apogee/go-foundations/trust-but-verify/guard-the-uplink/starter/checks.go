package main

import "testing"

// checkVerify checks a version of verify on a signed command and on forgeries of it.
func checkVerify(t *testing.T, verify func(cmd, sigHex string, key []byte) bool) {
	// TODO: one subtest each: valid, changed command, changed signature, wrong key, empty signature.
}
