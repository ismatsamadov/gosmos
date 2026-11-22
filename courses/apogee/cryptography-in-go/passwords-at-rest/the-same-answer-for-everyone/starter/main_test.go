package main

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"errors"
	"testing"
)

const (
	mirelaPassword = "fails closed, always"
	// mirelaRecord is Mirela's password at 600,000 iterations with the tests' readable salt.
	mirelaRecord = "pbkdf2-sha256$600000$6d6972656c6120746573742073616c74$6d91ceab779ff574b23820b133a740fa76d8b6a0deb5d6f0da2f94a6508dd938"
)

// testConsole returns a console holding Mirela's record, and the iteration count of every
// derivation its kdf has run so far.
func testConsole() (*Console, *[]int) {
	var counts []int
	c := &Console{
		records: map[string]string{"mirela": mirelaRecord},
		kdf: func(pw string, salt []byte, iter int) ([]byte, error) {
			counts = append(counts, iter)
			return pbkdf2.Key(sha256.New, pw, salt, iter, 32)
		},
	}
	return c, &counts
}

// TestSameAnswer tries a name nobody has, then a wrong password for a real one. Both must get
// errInvalidLogin after exactly one derivation at 600,000 iterations, so neither the answer
// nor the time it takes tells the two apart.
func TestSameAnswer(t *testing.T) {
	tests := []struct{ name, user, pw string }{
		{"an unknown user", "intruder", mirelaPassword},
		{"a wrong password", "mirela", "fails open, sometimes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, counts := testConsole()
			err := c.Check(tt.user, tt.pw)
			if !errors.Is(err, errInvalidLogin) || err.Error() != "invalid user or password" {
				t.Errorf("Check(%q, %q) = %v, want errInvalidLogin (invalid user or password)", tt.user, tt.pw, err)
			}
			if len(*counts) != 1 || (*counts)[0] != 600000 {
				t.Errorf("Check(%q, %q) ran kdf with the counts %v, want [600000]: one derivation, as a real record costs", tt.user, tt.pw, *counts)
			}
		})
	}
}

func TestRightPassword(t *testing.T) {
	c, _ := testConsole()
	if err := c.Check("mirela", mirelaPassword); err != nil {
		t.Errorf("Check(%q, %q) = %v, want nil", "mirela", mirelaPassword, err)
	}
}
