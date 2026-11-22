package main

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// The tests' passwords, and the test secrets in two engineers' authenticator apps: 20 ASCII
// bytes each. Sunniva's phone is in the snow at Ystmark, so Sunniva signs in from the sheet.
const (
	mirelaPassword  = "fails closed, always"
	sunnivaPassword = "wind 64, gusts 81"
	newPassword     = "two proofs, always"
	mirelasPhone    = "mirela's phone, 2028"
	sunnivasPhone   = "sunniva's phone 2028"
)

// clock is the console's clock in the tests: a set instant, moved by hand.
type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// world is the console on Friday 3 November 2028 at 16:00 UTC, as the tests see it: mirela with
// a phone and a sheet, sunniva with a sheet, and the kiosk's shared account with neither.
type world struct {
	gate    *Gate
	clock   *clock
	sheets  map[string][]string // user → their printed recovery codes
	emailed map[string]string   // user → the token in the last reset link queued for them
	saved   []string            // "<user> <record>" for each record c.save stored
}

// fastKDF stands in for PBKDF2 in these tests: 1,000 rounds whatever a record asks for, so a
// sign-in takes a millisecond instead of 60.
func fastKDF(pw string, salt []byte, iter int) ([]byte, error) {
	return pbkdf2.Key(sha256.New, pw, salt, 1000, 32)
}

func newWorld(t *testing.T) *world {
	w := &world{
		clock:   &clock{t: time.Date(2028, 11, 3, 16, 0, 0, 0, time.UTC)},
		sheets:  map[string][]string{},
		emailed: map[string]string{},
	}
	c := &Console{
		records: map[string]string{},
		kdf:     fastKDF,
		save: func(user, record string) error {
			w.saved = append(w.saved, user+" "+record)
			return nil
		},
		sessions: map[string]string{},
	}
	for user, pw := range map[string]string{"mirela": mirelaPassword, "sunniva": sunnivaPassword, "kiosk": "fishbowl"} {
		record, err := c.newRecord(pw)
		if err != nil {
			t.Fatalf("c.newRecord(%q) = %v", pw, err)
		}
		c.records[user] = record
	}
	recovery := &Recovery{unused: map[string][]string{}}
	w.sheets["mirela"] = recovery.Issue("mirela")
	w.sheets["sunniva"] = recovery.Issue("sunniva")
	w.gate = &Gate{
		console: c,
		codes: &Codes{
			secrets: map[string][]byte{"mirela": []byte(mirelasPhone), "sunniva": []byte(sunnivasPhone)},
			last:    map[string]uint64{},
			now:     w.clock.now,
		},
		recovery: recovery,
		resets: &Resets{
			pending: map[string]Reset{},
			now:     w.clock.now,
			known: func(user string) bool {
				_, ok := c.records[user]
				return ok
			},
			send: func(user, token string) { w.emailed[user] = token },
		},
		misses: map[string]int{},
	}
	return w
}

// code returns the code Mirela's phone shows now.
func (w *world) code() string {
	return totp([]byte(mirelasPhone), w.clock.t)
}

// wrongCode returns six digits that are none of the three codes Mirela's secret makes around now.
func (w *world) wrongCode() string {
	step := timeStep(w.clock.t)
	right := map[string]bool{}
	for s := step - 1; s <= step+1; s++ {
		right[hotp([]byte(mirelasPhone), s)] = true
	}
	for n := 0; ; n++ {
		if code := fmt.Sprintf("%06d", n); !right[code] {
			return code
		}
	}
}

// wait moves the clock on 30 seconds, to the next step and its new code.
func (w *world) wait() {
	w.clock.t = w.clock.t.Add(30 * time.Second)
}

// refused reports whether a sign-in was refused with the console's one answer and no token.
func refused(token string, err error) bool {
	return token == "" && errors.Is(err, errInvalidSignIn)
}

// TestPasswordAlone plays a thief who has Mirela's password and not her phone.
func TestPasswordAlone(t *testing.T) {
	w := newWorld(t)
	for _, code := range []string{"", w.wrongCode(), w.sheets["sunniva"][0]} {
		if token, err := w.gate.SignIn("mirela", mirelaPassword, code); !refused(token, err) {
			t.Errorf("SignIn(%q, the right password, %q) = %q, %v, want \"\", errInvalidSignIn", "mirela", code, token, err)
		}
	}
	if n := len(w.gate.console.sessions); n != 0 {
		t.Errorf("three refused sign-ins left %d sessions, want 0", n)
	}
}

// TestPasswordAndCode signs Mirela in with the code her phone shows, then tries that code again.
func TestPasswordAndCode(t *testing.T) {
	w := newWorld(t)
	code := w.code()
	token, err := w.gate.SignIn("mirela", mirelaPassword, code)
	if err != nil || len(token) < 26 {
		t.Fatalf("SignIn(%q, the right password, %q) = %q, %v, want a session token, nil", "mirela", code, token, err)
	}
	if user := w.gate.console.sessions[tokenDigest(token)]; user != "mirela" {
		t.Errorf("c.sessions[tokenDigest(token)] = %q, want %q", user, "mirela")
	}
	w.clock.t = w.clock.t.Add(10 * time.Second)
	if token, err := w.gate.SignIn("mirela", mirelaPassword, code); !refused(token, err) {
		t.Errorf("SignIn(%q, the right password, %q) again 10 seconds later = %q, %v, want \"\", errInvalidSignIn", "mirela", code, token, err)
	}
	w.wait()
	if _, err := w.gate.SignIn("mirela", mirelaPassword, w.code()); err != nil {
		t.Errorf("SignIn(%q, the right password, the next step's code) = %v, want nil", "mirela", err)
	}
}

// TestRecoveryCodes signs Sunniva in from the sheet, each code once.
func TestRecoveryCodes(t *testing.T) {
	w := newWorld(t)
	sheet := w.sheets["sunniva"]
	if _, err := w.gate.SignIn("sunniva", sunnivaPassword, sheet[0]); err != nil {
		t.Fatalf("SignIn(%q, the right password, recovery code 1) = %v, want nil", "sunniva", err)
	}
	if token, err := w.gate.SignIn("sunniva", sunnivaPassword, sheet[0]); !refused(token, err) {
		t.Errorf("SignIn(%q, the right password, recovery code 1) again = %q, %v, want \"\", errInvalidSignIn", "sunniva", token, err)
	}
	if _, err := w.gate.SignIn("sunniva", sunnivaPassword, sheet[1]); err != nil {
		t.Errorf("SignIn(%q, the right password, recovery code 2) = %v, want nil", "sunniva", err)
	}
}

// TestNoPhoneNoSheet tries the kiosk's shared account, which has neither: nothing signs it in,
// not even the codes anyone can make from an empty secret.
func TestNoPhoneNoSheet(t *testing.T) {
	w := newWorld(t)
	step := timeStep(w.clock.t)
	for _, code := range []string{"", hotp(nil, step-1), hotp(nil, step), hotp(nil, step+1)} {
		if token, err := w.gate.SignIn("kiosk", "fishbowl", code); !refused(token, err) {
			t.Errorf("SignIn(%q, the right password, %q) = %q, %v, want \"\", errInvalidSignIn", "kiosk", code, token, err)
		}
	}
	if token, err := w.gate.SignIn("intruder", "fishbowl", w.code()); !refused(token, err) {
		t.Errorf("SignIn(%q, %q, a code) = %q, %v, want \"\", errInvalidSignIn", "intruder", "fishbowl", token, err)
	}
}

// TestWrongPasswordSpendsNothing sends good codes with wrong passwords first: the codes must
// still work with the right ones.
func TestWrongPasswordSpendsNothing(t *testing.T) {
	w := newWorld(t)
	code, sheet := w.code(), w.sheets["sunniva"]
	if token, err := w.gate.SignIn("mirela", "fails open", code); !refused(token, err) {
		t.Errorf("SignIn(%q, a wrong password, the phone's code) = %q, %v, want \"\", errInvalidSignIn", "mirela", token, err)
	}
	if token, err := w.gate.SignIn("sunniva", "wind 72", sheet[0]); !refused(token, err) {
		t.Errorf("SignIn(%q, a wrong password, recovery code 1) = %q, %v, want \"\", errInvalidSignIn", "sunniva", token, err)
	}
	if _, err := w.gate.SignIn("mirela", mirelaPassword, code); err != nil {
		t.Errorf("SignIn(%q, the right password, the code a wrong password came with) = %v, want nil: a refused password spends nothing", "mirela", err)
	}
	if _, err := w.gate.SignIn("sunniva", sunnivaPassword, sheet[0]); err != nil {
		t.Errorf("SignIn(%q, the right password, recovery code 1 again) = %v, want nil: a refused password spends nothing", "sunniva", err)
	}
}

// TestTenWrongCodes checks the count of wrong codes in a row, on either side of maxMisses.
func TestTenWrongCodes(t *testing.T) {
	miss := func(w *world, n int) {
		for range n {
			w.gate.SignIn("mirela", mirelaPassword, w.wrongCode())
		}
	}
	t.Run("nine wrong codes, then the right one", func(t *testing.T) {
		w := newWorld(t)
		miss(w, 9)
		if _, err := w.gate.SignIn("mirela", mirelaPassword, w.code()); err != nil {
			t.Errorf("SignIn(%q, the right password, the right code) after 9 wrong codes = %v, want nil", "mirela", err)
		}
	})
	t.Run("ten wrong codes", func(t *testing.T) {
		w := newWorld(t)
		miss(w, 10)
		w.wait()
		if token, err := w.gate.SignIn("mirela", mirelaPassword, w.code()); !refused(token, err) {
			t.Errorf("SignIn(%q, the right password, the right code) after 10 wrong codes = %q, %v, want \"\", errInvalidSignIn", "mirela", token, err)
		}
		if token, err := w.gate.SignIn("mirela", mirelaPassword, w.sheets["mirela"][0]); !refused(token, err) {
			t.Errorf("SignIn(%q, the right password, recovery code 1) after 10 wrong codes = %q, %v, want \"\", errInvalidSignIn", "mirela", token, err)
		}
	})
	t.Run("a right code starts the count again", func(t *testing.T) {
		w := newWorld(t)
		miss(w, 9)
		if _, err := w.gate.SignIn("mirela", mirelaPassword, w.code()); err != nil {
			t.Fatalf("SignIn(%q, the right password, the right code) after 9 wrong codes = %v, want nil", "mirela", err)
		}
		miss(w, 9)
		w.wait()
		if _, err := w.gate.SignIn("mirela", mirelaPassword, w.code()); err != nil {
			t.Errorf("SignIn(%q, the right password, the right code) after 9, 1 right and 9 more wrong codes = %v, want nil", "mirela", err)
		}
	})
	t.Run("wrong passwords count nothing", func(t *testing.T) {
		w := newWorld(t)
		for range 20 {
			w.gate.SignIn("mirela", "fails open", w.wrongCode())
		}
		if _, err := w.gate.SignIn("mirela", mirelaPassword, w.code()); err != nil {
			t.Errorf("SignIn(%q, the right password, the right code) after 20 wrong passwords = %v, want nil: only a right password lets a code count", "mirela", err)
		}
	})
}

// TestReset walks Mirela's forgotten password through: a link, a new password, and still a code.
func TestReset(t *testing.T) {
	w := newWorld(t)
	before, err := w.gate.SignIn("mirela", mirelaPassword, w.code())
	if err != nil {
		t.Fatalf("SignIn(%q, the right password, the right code) = %v, want nil", "mirela", err)
	}
	w.gate.resets.Request("mirela")
	token := w.emailed["mirela"]
	w.saved = nil
	if err := w.gate.ResetPassword("mirela", token, newPassword); err != nil {
		t.Fatalf("ResetPassword(%q, the emailed token, %q) = %v, want nil", "mirela", newPassword, err)
	}
	if len(w.saved) != 1 || !strings.HasPrefix(w.saved[0], "mirela pbkdf2-sha256$600000$") {
		t.Errorf("ResetPassword saved %d records, want 1 new record for mirela through c.save", len(w.saved))
	}
	if _, open := w.gate.console.sessions[tokenDigest(before)]; open {
		t.Errorf("Mirela's session from before the reset is still in c.sessions, want it ended")
	}
	w.wait()
	if token, err := w.gate.SignIn("mirela", mirelaPassword, w.code()); !refused(token, err) {
		t.Errorf("SignIn(%q, the old password, the right code) after the reset = %q, %v, want \"\", errInvalidSignIn", "mirela", token, err)
	}
	if token, err := w.gate.SignIn("mirela", newPassword, ""); !refused(token, err) {
		t.Errorf("SignIn(%q, the new password, no code) = %q, %v, want \"\", errInvalidSignIn", "mirela", token, err)
	}
	if _, err := w.gate.SignIn("mirela", newPassword, w.code()); err != nil {
		t.Errorf("SignIn(%q, the new password, the right code) = %v, want nil", "mirela", err)
	}
	if err := w.gate.ResetPassword("mirela", token, "a third one"); !errors.Is(err, errInvalidReset) {
		t.Errorf("ResetPassword(%q, the same token, %q) = %v, want errInvalidReset", "mirela", "a third one", err)
	}
}

// TestResetUnlocks locks Mirela's account with ten wrong codes: only a reset opens it again,
// and the recovery code tried while it was locked still works afterwards.
func TestResetUnlocks(t *testing.T) {
	w := newWorld(t)
	for range 10 {
		w.gate.SignIn("mirela", mirelaPassword, w.wrongCode())
	}
	sheet := w.sheets["mirela"]
	if token, err := w.gate.SignIn("mirela", mirelaPassword, sheet[0]); !refused(token, err) {
		t.Fatalf("SignIn(%q, the right password, recovery code 1) after 10 wrong codes = %q, %v, want \"\", errInvalidSignIn", "mirela", token, err)
	}
	w.gate.resets.Request("mirela")
	if err := w.gate.ResetPassword("mirela", w.emailed["mirela"], newPassword); err != nil {
		t.Fatalf("ResetPassword(%q, the emailed token, %q) = %v, want nil", "mirela", newPassword, err)
	}
	if _, err := w.gate.SignIn("mirela", newPassword, sheet[0]); err != nil {
		t.Errorf("SignIn(%q, the new password, recovery code 1) after the reset = %v, want nil: a locked account spends no code", "mirela", err)
	}
}

// TestBadResetChangesNothing sends reset requests that must fail, and checks that each one
// leaves Mirela's password, her session and the user file as they were.
func TestBadResetChangesNothing(t *testing.T) {
	w := newWorld(t)
	session, err := w.gate.SignIn("mirela", mirelaPassword, w.code())
	if err != nil {
		t.Fatalf("SignIn(%q, the right password, the right code) = %v, want nil", "mirela", err)
	}
	w.gate.resets.Request("mirela")
	token := w.emailed["mirela"]
	w.saved = nil
	tests := []struct{ name, user, token string }{
		{"no token", "mirela", ""},
		{"the digest the console keeps", "mirela", tokenDigest(token)},
		{"another account's name", "sunniva", token},
	}
	for _, tt := range tests {
		if err := w.gate.ResetPassword(tt.user, tt.token, "stolen"); !errors.Is(err, errInvalidReset) {
			t.Errorf("ResetPassword(%q, %s, %q) = %v, want errInvalidReset", tt.user, tt.name, "stolen", err)
		}
	}
	if len(w.saved) != 0 {
		t.Errorf("refused resets saved %d records, want none", len(w.saved))
	}
	if _, open := w.gate.console.sessions[tokenDigest(session)]; !open {
		t.Errorf("Mirela's session is gone after refused resets, want it still open")
	}
	w.wait()
	if _, err := w.gate.SignIn("mirela", mirelaPassword, w.code()); err != nil {
		t.Errorf("SignIn(%q, the old password, the right code) after refused resets = %v, want nil", "mirela", err)
	}
}
