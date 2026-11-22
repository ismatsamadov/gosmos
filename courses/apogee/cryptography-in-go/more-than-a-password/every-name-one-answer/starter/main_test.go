package main

import (
	"strings"
	"testing"
	"time"
)

// clock is the console's clock in the tests: a set instant, moved by hand.
type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// mail holds what the tests' reset form queued: "<user> <token>" for each email, in order.
type mail struct{ sent []string }

// testResets returns an empty reset store on Thursday 26 October 2028 at 14:00 UTC, for a
// console whose accounts are mirela, nkem, ingrid and kiosk, and the mail it queues.
func testResets() (*Resets, *mail) {
	accounts := map[string]bool{"mirela": true, "nkem": true, "ingrid": true, "kiosk": true}
	m := &mail{}
	c := &clock{t: time.Date(2028, 10, 26, 14, 0, 0, 0, time.UTC)}
	r := &Resets{
		pending: map[string]Reset{},
		now:     c.now,
		known:   func(user string) bool { return accounts[user] },
		send:    func(user, token string) { m.sent = append(m.sent, user+" "+token) },
	}
	return r, m
}

// TestFormListsNoOne plays Mirela at the reset form: names with an account and names without
// one must get the same answer, or the form reads out the staff list.
func TestFormListsNoOne(t *testing.T) {
	r, _ := testResets()
	first := r.Request("mirela")
	for _, user := range []string{"intruder", "nkem", "admin", "ingrid", "root", "kiosk", "paz.o", ""} {
		if got := r.Request(user); got != first {
			t.Errorf("Request(%q) = %q, but Request(%q) = %q, want one answer for every name", user, got, "mirela", first)
		}
	}
	if first != resetSent {
		t.Errorf("Request(%q) = %q, want resetSent, %q", "mirela", first, resetSent)
	}
}

// TestMailOnlyToAccounts checks who gets a link: the owner of an account that exists, and
// nobody else.
func TestMailOnlyToAccounts(t *testing.T) {
	r, m := testResets()
	for _, user := range []string{"intruder", "root"} {
		r.Request(user)
	}
	if len(m.sent) != 0 || len(r.pending) != 0 {
		t.Fatalf("Request for two names without accounts queued %d emails and kept %d resets, want none", len(m.sent), len(r.pending))
	}
	r.Request("nkem")
	if len(m.sent) != 1 {
		t.Fatalf("Request(%q) queued %d emails, want 1", "nkem", len(m.sent))
	}
	user, token, _ := strings.Cut(m.sent[0], " ")
	if user != "nkem" {
		t.Errorf("Request(%q) emailed %q, want %q", "nkem", user, "nkem")
	}
	if err := r.Redeem("nkem", token); err != nil {
		t.Errorf("Redeem(%q, the emailed token) = %v, want nil", "nkem", err)
	}
}
