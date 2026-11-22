package main

import "testing"

// checkVerify checks a version of verify on a signed command and on forgeries of it.
func checkVerify(t *testing.T, verify func(cmd, sigHex string, key []byte) bool) {
	sig := sign(labCommand, []byte(labKey))
	cases := []struct {
		name, cmd, sig, key string
		want                bool
	}{
		{"valid", labCommand, sig, labKey, true},
		{"changed command", "pipit-2|118|SET_MODE|NOMINAL", sig, labKey, false},
		{"changed signature", labCommand, sig[:63] + "0", labKey, false},
		{"wrong key", labCommand, sig, "someone else's key", false},
		{"empty signature", labCommand, "", labKey, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := verify(c.cmd, c.sig, []byte(c.key)); got != c.want {
				t.Errorf("verify(%q, %q, %q) = %v, want %v", c.cmd, c.sig, c.key, got, c.want)
			}
		})
	}
}
