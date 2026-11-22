package main

import (
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

// loaderSettings are a dish loader's settings; only the loader reads the
// password of its host's line.
type loaderSettings struct {
	Station  string
	Queue    int
	Uploads  int
	password string
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

// sameLines fails the test at the first line where got and want differ.
func sameLines(t *testing.T, call string, got, want []string) {
	t.Helper()
	for i := range max(len(got), len(want)) {
		switch {
		case len(want) == 0:
			t.Errorf("%s = %q, want nil", call, got)
		case i == len(got):
			t.Errorf("%s gave %d of %d lines: %q is missing", call, len(got), len(want), want[i])
		case i == len(want):
			t.Errorf("%s line %d = %q, want none after line %d", call, i+1, got[i], i)
		case got[i] != want[i]:
			t.Errorf("%s line %d = %q, want %q", call, i+1, got[i], want[i])
		default:
			continue
		}
		return
	}
}

func TestDescribe(t *testing.T) {
	all := []string{"Slots=12", "Days=100", "MinElevDeg=10", "Turnaround=2m0s", "Replan=10m0s", "RunLimit=30s", "Queue=256", "Quiet=5s"}
	cases := []struct {
		name, call string
		v          any
		want       []string
	}{
		{"passplan 2's settings", "describe(defaults)", defaults, all},
		{"the same settings through a pointer", "describe(&defaults)", &defaults, all},
		{"Vorrel Island's dish loader, password and all", `describe(loaderSettings{"VRL", 256, 2, "…"})`,
			loaderSettings{Station: "VRL", Queue: 256, Uploads: 2, password: "ice-4-ever"},
			[]string{"Station=VRL", "Queue=256", "Uploads=2"}},
		{"a slot count, not a struct", "describe(12)", 12, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := safely(t, c.call, func() []string { return describe(c.v) })
			sameLines(t, c.call, got, c.want)
		})
	}
}
