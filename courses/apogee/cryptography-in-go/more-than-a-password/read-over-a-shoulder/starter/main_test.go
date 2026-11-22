package main

import (
	"testing"
	"time"
)

// The test secrets in Nkem's and Mirela's authenticator apps: 20 ASCII bytes each.
const (
	nkemsPhone   = "nkem's phone in 2028"
	mirelasPhone = "mirela's phone, 2028"
)

// clock is the console's clock in the tests: a set instant, moved by hand.
type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// at returns a time on Wednesday 1 November 2028, UTC; 10:15:00 starts a step.
func at(hour, min, sec int) time.Time {
	return time.Date(2028, 11, 1, hour, min, sec, 0, time.UTC)
}

// testCodes returns a code checker for nkem and mirela, nothing used yet, and its clock, set to at.
func testCodes(at time.Time) (*Codes, *clock) {
	c := &clock{t: at}
	return &Codes{
		secrets: map[string][]byte{"nkem": []byte(nkemsPhone), "mirela": []byte(mirelasPhone)},
		last:    map[string]uint64{},
		now:     c.now,
	}, c
}

// TestReadOverAShoulder plays Mirela: she reads the code Nkem types at 10:15:05 and types it
// herself 40 seconds later, while the code is still inside its window.
func TestReadOverAShoulder(t *testing.T) {
	codes, clock := testCodes(at(10, 15, 5))
	code := totp([]byte(nkemsPhone), at(10, 15, 5))
	if !codes.Check("nkem", code) {
		t.Fatalf("Check(%q, %q) at 10:15:05 = false, want true", "nkem", code)
	}
	clock.t = at(10, 15, 45)
	if codes.Check("nkem", code) {
		t.Errorf("Check(%q, %q) again at 10:15:45 = true, want false: a code counts once", "nkem", code)
	}
	clock.t = at(10, 15, 50)
	if next := totp([]byte(nkemsPhone), at(10, 15, 50)); !codes.Check("nkem", next) {
		t.Errorf("Check(%q, %q), the next step's code, at 10:15:50 = false, want true", "nkem", next)
	}
}

// TestFastPhone checks a phone 20 seconds ahead: its code is the next step's, and once that
// step is used, neither it nor an older step's code counts again.
func TestFastPhone(t *testing.T) {
	codes, clock := testCodes(at(10, 15, 20))
	ahead := totp([]byte(nkemsPhone), at(10, 15, 40))
	if !codes.Check("nkem", ahead) {
		t.Fatalf("Check(%q, %q), a code from 20 seconds ahead, at 10:15:20 = false, want true", "nkem", ahead)
	}
	clock.t = at(10, 15, 25)
	if codes.Check("nkem", ahead) {
		t.Errorf("Check(%q, %q), the same code, at 10:15:25 = true, want false", "nkem", ahead)
	}
	if older := totp([]byte(nkemsPhone), at(10, 15, 25)); codes.Check("nkem", older) {
		t.Errorf("Check(%q, %q), the step before the one used, at 10:15:25 = true, want false", "nkem", older)
	}
}

// TestRefusalChangesNothing sends a wrong code first: it must leave nothing behind.
func TestRefusalChangesNothing(t *testing.T) {
	codes, _ := testCodes(at(10, 15, 5))
	right := totp([]byte(nkemsPhone), at(10, 15, 5))
	wrong := totp([]byte(mirelasPhone), at(10, 15, 5))
	if codes.Check("nkem", wrong) {
		t.Fatalf("Check(%q, %q), Mirela's code, = true, want false", "nkem", wrong)
	}
	if last := codes.last["nkem"]; last != 0 {
		t.Errorf("codes.last[%q] = %d after a refused code, want 0: a refusal changes nothing", "nkem", last)
	}
	if !codes.Check("nkem", right) {
		t.Errorf("Check(%q, %q) after a refused code = false, want true", "nkem", right)
	}
}

// TestEachAccountApart checks that one engineer's code doesn't use up another's step.
func TestEachAccountApart(t *testing.T) {
	codes, _ := testCodes(at(10, 15, 5))
	for _, user := range []string{"mirela", "nkem"} {
		secret := map[string]string{"mirela": mirelasPhone, "nkem": nkemsPhone}[user]
		if code := totp([]byte(secret), at(10, 15, 5)); !codes.Check(user, code) {
			t.Errorf("Check(%q, %q) at 10:15:05 = false, want true", user, code)
		}
	}
	if code := hotp(nil, timeStep(at(10, 15, 5))); codes.Check("kiosk", code) {
		t.Errorf("Check(%q, %q), the code an empty secret makes, = true, want false", "kiosk", code)
	}
}
