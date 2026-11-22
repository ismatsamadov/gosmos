package main

import (
	"errors"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// firstFires is a Debouncer that also fires at once, on the first Notify.
type firstFires struct {
	*Debouncer
	once sync.Once
	fire func()
}

// Notify fires on the first call, then debounces as the Debouncer does.
func (m *firstFires) Notify() {
	m.once.Do(m.fire)
	m.Debouncer.Notify()
}

// neverResets is a Debouncer that ignores every Notify after the first.
type neverResets struct {
	*Debouncer
	once sync.Once
}

// Notify starts the quiet on the first call only.
func (m *neverResets) Notify() { m.once.Do(m.Debouncer.Notify) }

// ignoresStop is a Debouncer whose Stop does nothing.
type ignoresStop struct{ *Debouncer }

// Stop does nothing.
func (ignoresStop) Stop() {}

// mutants are Nkem's broken copies of the debouncer, each with one bug planted.
var mutants = []struct {
	bug  string
	newD func(quiet time.Duration, fire func()) notifier
}{
	{"fires at the first Notify", func(quiet time.Duration, fire func()) notifier {
		return &firstFires{Debouncer: newDebouncer(quiet, fire), fire: fire}
	}},
	{"never resets", func(quiet time.Duration, fire func()) notifier {
		return &neverResets{Debouncer: newDebouncer(quiet, fire)}
	}},
	{"fires twice", func(quiet time.Duration, fire func()) notifier {
		return newDebouncer(quiet, func() { fire(); fire() })
	}},
	{"ignores Stop", func(quiet time.Duration, fire func()) notifier {
		return ignoresStop{newDebouncer(quiet, fire)}
	}},
}

// debouncer returns the debouncer this run checks: yours, or the mutant APOGEE_MUTANT names.
func debouncer(t *testing.T) func(quiet time.Duration, fire func()) notifier {
	bug := os.Getenv("APOGEE_MUTANT")
	if bug == "" {
		return func(quiet time.Duration, fire func()) notifier { return newDebouncer(quiet, fire) }
	}
	for _, m := range mutants {
		if m.bug == bug {
			return m.newD
		}
	}
	t.Fatalf("no mutant %q", bug)
	return nil
}

func TestCheckDebouncer(t *testing.T) {
	newD := debouncer(t)
	start := time.Now()             // the real clock; a bubble's fake one starts in 2000
	var made, realClock atomic.Bool // realClock: a debouncer made on the real clock, outside a bubble
	// Twenty runs: a check that looks without synctest.Wait fails some of them.
	for range 20 {
		checkDebouncer(t, func(quiet time.Duration, fire func()) notifier {
			made.Store(true)
			if d := time.Since(start); d >= 0 && d < time.Minute {
				realClock.Store(true)
			}
			return newD(quiet, fire)
		})
		if t.Failed() || realClock.Load() {
			break
		}
	}
	if os.Getenv("APOGEE_MUTANT") != "" {
		return
	}
	if !made.Load() {
		t.Errorf("checkDebouncer never called newD, want a debouncer made and checked")
	}
	if realClock.Load() {
		t.Errorf("checkDebouncer made a debouncer on the real clock, want every one inside synctest.Test")
	}
}

func TestChecksCatchMutants(t *testing.T) {
	for _, m := range mutants {
		t.Run("the mutant that "+m.bug, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestCheckDebouncer$", "-test.timeout=5s")
			cmd.Env = append(os.Environ(), "APOGEE_MUTANT="+m.bug)
			err := cmd.Run()
			var exit *exec.ExitError
			switch {
			case err == nil:
				t.Errorf("the mutant that %s survived: checkDebouncer reported nothing", m.bug)
			case !errors.As(err, &exit):
				t.Fatalf("running the mutant that %s: %v", m.bug, err)
			}
		})
	}
}
