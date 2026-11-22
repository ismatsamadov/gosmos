package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// defaults are passplan 2's settings as their default tags give them.
var defaults = Config{
	Slots:      12,
	Days:       100,
	MinElevDeg: 10,
	Turnaround: 2 * time.Minute,
	Replan:     10 * time.Minute,
	RunLimit:   30 * time.Second,
	Queue:      256,
	Quiet:      5 * time.Second,
}

// deploy is what a test deploy of passplan 2 reads when it starts.
type deploy struct {
	Label string `check:"required"`
	Plan  Config
	token string `check:"required"` // only the deploy tool reads it
	last  Config // the deploy before, kept for the deploy tool
}

// loaderSetup is what a new station's dish loader reads: its host's
// form, then its own settings.
type loaderSetup struct {
	HostConfig
	Uploads int    `check:"min=1,max=2"`
	Remark  string `check:""` // rules to come
}

// safely calls f and fails the test, naming call, if f panics.
func safely[T any](t *testing.T, call string, f func() T) T {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("%s panicked: %v", call, r)
		}
	}()
	return f()
}

// sameProblems fails the test at the first problem in err that differs
// from want; a nil err has none.
func sameProblems(t *testing.T, call string, err error, want []string) {
	t.Helper()
	var got []string
	if err != nil {
		got = strings.Split(err.Error(), "\n")
	}
	for i := range max(len(got), len(want)) {
		switch {
		case len(want) == 0:
			t.Errorf("%s = %q, want nil", call, err)
		case len(got) == 0 && len(want) == 1:
			t.Errorf("%s = nil, want %q", call, want[0])
		case len(got) == 0:
			t.Errorf("%s = nil, want %d problems, the first %q", call, len(want), want[0])
		case i == len(got):
			t.Errorf("%s gave %d of %d problems: %q is missing", call, len(got), len(want), want[i])
		case i == len(want):
			t.Errorf("%s problem %d = %q, want none after problem %d", call, i+1, got[i], i)
		case got[i] != want[i]:
			t.Errorf("%s problem %d = %q, want %q", call, i+1, got[i], want[i])
		default:
			continue
		}
		return
	}
}

func TestValidate(t *testing.T) {
	t.Run("the fifteen forms as the hosts first sent them", func(t *testing.T) {
		want := map[string][]string{
			"ICM": {"HostConfig.Name: required", `HostConfig.Bands: "K" is not one of UHF S X`},
			"VRL": {"HostConfig.Dishes: 0 is below min 1"},
		}
		for i, f := range hostForms {
			t.Run(f.Code, func(t *testing.T) {
				call := fmt.Sprintf("Validate(hostForms[%d])", i)
				sameProblems(t, call, safely(t, call, func() error { return Validate(f) }), want[f.Code])
			})
		}
	})
	t.Run("passplan 2's defaults", func(t *testing.T) {
		sameProblems(t, "Validate(defaults)", safely(t, "Validate(defaults)", func() error { return Validate(defaults) }), nil)
	})
	t.Run("an empty Config", func(t *testing.T) {
		err := safely(t, "Validate(Config{})", func() error { return Validate(Config{}) })
		sameProblems(t, "Validate(Config{})", err, []string{
			"Config.Slots: 0 is below min 1",
			"Config.Days: 0 is below min 1",
			"Config.Replan: 0s is below min 1m0s",
			"Config.RunLimit: 0s is below min 1s",
			"Config.Queue: 0 is below min 1",
		})
		if _, joined := err.(interface{ Unwrap() []error }); err != nil && !joined {
			t.Errorf("Validate(Config{}) is one error, want the problems joined with errors.Join")
		}
	})
	t.Run("an empty Config after applyDefaults", func(t *testing.T) {
		var c Config
		if err := applyDefaults(&c); err != nil {
			t.Fatalf("applyDefaults(&c) = %v", err)
		}
		sameProblems(t, "Validate(c)", safely(t, "Validate(c)", func() error { return Validate(c) }), nil)
	})
	t.Run("a pointer to settings with two bad fields", func(t *testing.T) {
		c := defaults
		c.Slots, c.Turnaround = 70, 15*time.Minute
		err := safely(t, "Validate(&c)", func() error { return Validate(&c) })
		sameProblems(t, "Validate(&c)", err, []string{
			"Config.Slots: 70 is above max 64",
			"Config.Turnaround: 15m0s is above max 10m0s",
		})
	})
	t.Run("Friday's test deploy", func(t *testing.T) {
		d := deploy{Label: "test-0601", Plan: defaults}
		d.Plan.Slots = 0
		call := `Validate(deploy{Label: "test-0601", Plan: …})`
		sameProblems(t, call, safely(t, call, func() error { return Validate(d) }), []string{"Config.Slots: 0 is below min 1"})
	})
	t.Run("Vorrel Island's dish loader", func(t *testing.T) {
		l := loaderSetup{HostConfig: hostForms[14], Uploads: 3}
		call := "Validate(loaderSetup{HostConfig: hostForms[14], Uploads: 3})"
		sameProblems(t, call, safely(t, call, func() error { return Validate(l) }), []string{
			"HostConfig.Dishes: 0 is below min 1",
			"loaderSetup.Uploads: 3 is above max 2",
		})
	})
	t.Run("a slot count, not a struct", func(t *testing.T) {
		sameProblems(t, "Validate(12)", safely(t, "Validate(12)", func() error { return Validate(12) }), []string{"validate: want a struct, got int"})
	})
	t.Run("all fifteen forms at once, a slice", func(t *testing.T) {
		err := safely(t, "Validate(hostForms)", func() error { return Validate(hostForms) })
		sameProblems(t, "Validate(hostForms)", err, []string{"validate: want a struct, got []main.HostConfig"})
	})
}
