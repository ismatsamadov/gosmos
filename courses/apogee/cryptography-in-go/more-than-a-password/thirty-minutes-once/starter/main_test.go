package main

import (
	"errors"
	"testing"
	"time"
)

// clock is the console's clock in the tests: a set instant, moved by hand.
type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// started is when the tests' resets start: Thursday 26 October 2028, 09:00 UTC.
var started = time.Date(2028, 10, 26, 9, 0, 0, 0, time.UTC)

// testResets returns an empty reset store and the clock it reads, set to started.
func testResets() (*Resets, *clock) {
	c := &clock{t: started}
	return &Resets{pending: map[string]Reset{}, now: c.now}, c
}

// refused reports whether err is the console's one answer to a link it won't take, unwrapped,
// so that no refusal says which check failed.
func refused(err error) bool {
	return errors.Is(err, errInvalidReset) && err.Error() == "invalid or expired reset link"
}

// short cuts a value longer than 50 characters to its first 8 and an ellipsis.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

func TestExpires(t *testing.T) {
	tests := []struct {
		name  string
		after time.Duration
		ok    bool
	}{
		{"at once", 0, true},
		{"at 09:29:59", 30*time.Minute - time.Second, true},
		{"at 09:30:00 exactly", 30 * time.Minute, false},
		{"the next day", 24 * time.Hour, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, c := testResets()
			token := r.Start("mirela")
			c.t = started.Add(tt.after)
			err := r.Redeem("mirela", token)
			if tt.ok && err != nil {
				t.Errorf("Redeem(%q, the link's token) %s = %v, want nil", "mirela", tt.name, err)
			}
			if !tt.ok && !refused(err) {
				t.Errorf("Redeem(%q, the link's token) %s = %v, want errInvalidReset", "mirela", tt.name, err)
			}
		})
	}
}

// TestOnce uses one link twice, as whoever finds the old email could.
func TestOnce(t *testing.T) {
	r, c := testResets()
	token := r.Start("nkem")
	c.t = started.Add(5 * time.Minute)
	if err := r.Redeem("nkem", token); err != nil {
		t.Fatalf("Redeem(%q, the link's token) = %v, want nil", "nkem", err)
	}
	c.t = started.Add(6 * time.Minute)
	if err := r.Redeem("nkem", token); !refused(err) {
		t.Errorf("Redeem(%q, the link's token) a second time = %v, want errInvalidReset", "nkem", err)
	}
	if err := r.Redeem("intruder", token); !refused(err) {
		t.Errorf("Redeem(%q, nkem's token) = %v, want errInvalidReset", "intruder", err)
	}
}

// TestRefusalChangesNothing sends wrong tokens first: none of them may spend the real link.
func TestRefusalChangesNothing(t *testing.T) {
	r, c := testResets()
	token := r.Start("ingrid")
	for _, wrong := range []string{"", "AAAAAAAAAAAAAAAAAAAAAAAAAA", tokenDigest(token)} {
		if err := r.Redeem("ingrid", wrong); !refused(err) {
			t.Errorf("Redeem(%q, %q) = %v, want errInvalidReset", "ingrid", short(wrong), err)
		}
	}
	c.t = started.Add(10 * time.Minute)
	if err := r.Redeem("ingrid", token); err != nil {
		t.Errorf("Redeem(%q, the link's token) after three wrong tokens = %v, want nil: a refusal changes nothing", "ingrid", err)
	}
}

// TestNewestLinkOnly starts two resets for one account: only the second link works.
func TestNewestLinkOnly(t *testing.T) {
	r, _ := testResets()
	first := r.Start("paz")
	second := r.Start("paz")
	if err := r.Redeem("paz", first); !refused(err) {
		t.Errorf("Redeem(%q, the first link's token) after a second Start = %v, want errInvalidReset", "paz", err)
	}
	if err := r.Redeem("paz", second); err != nil {
		t.Errorf("Redeem(%q, the second link's token) = %v, want nil", "paz", err)
	}
}
