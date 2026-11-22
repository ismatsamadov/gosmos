package main

import (
	"fmt"
	"reflect"
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

// loader is a dish loader's settings; only the loader itself sets the
// password of its host's line.
type loader struct {
	Station  string
	Queue    int    `default:"256"`
	Uploads  int    `default:"2"`
	password string `default:"changeme"`
}

// typo is a settings struct whose default was written in words.
type typo struct {
	Slots int `default:"twelve"`
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

// firstDiff names the first field in which got and want, two values
// of one struct type, differ, or returns "" when none does.
func firstDiff(got, want any) string {
	g, w := reflect.ValueOf(got), reflect.ValueOf(want)
	for i := range g.NumField() {
		gs, ws := fmt.Sprint(g.Field(i)), fmt.Sprint(w.Field(i))
		if gs != ws {
			return fmt.Sprintf("%s = %s, want %s", g.Type().Field(i).Name, gs, ws)
		}
	}
	return ""
}

func TestApplyDefaults(t *testing.T) {
	six := defaults
	six.Slots = 6
	n := 12
	cases := []struct {
		name, call string
		ptr        any
		want       any    // the struct ptr points to afterward; nil when ptr isn't a pointer to a struct
		err        string // the error's text; "" for none
	}{
		{"an empty Config", "applyDefaults(&Config{})", &Config{}, defaults, ""},
		{"six slots, set on purpose", "applyDefaults(&Config{Slots: 6})", &Config{Slots: 6}, six, ""},
		{"tlmd's settings, with no default for the feeds", "applyDefaults(&PipelineConfig{})",
			&PipelineConfig{}, PipelineConfig{Workers: 4, Buffer: 16}, ""},
		{"a dish loader with a password default", `applyDefaults(&loader{Station: "VRL"})`,
			&loader{Station: "VRL"}, loader{Station: "VRL", Queue: 256, Uploads: 2}, ""},
		{"a Config, not a pointer to one", "applyDefaults(Config{})", Config{}, nil, "defaults: want a pointer to a struct, got main.Config"},
		{"a pointer to a slot count", "applyDefaults(&n)", &n, nil, "defaults: want a pointer to a struct, got *int"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := safely(t, c.call, func() error { return applyDefaults(c.ptr) })
			switch {
			case c.err == "" && err != nil:
				t.Fatalf("%s error = %q, want nil", c.call, err)
			case c.err != "" && err == nil:
				t.Fatalf("%s error = nil, want %q", c.call, c.err)
			case c.err != "" && err.Error() != c.err:
				t.Fatalf("%s error = %q, want %q", c.call, err, c.err)
			}
			if c.want == nil {
				return
			}
			if d := firstDiff(reflect.ValueOf(c.ptr).Elem().Interface(), c.want); d != "" {
				t.Errorf("after %s, %s", c.call, d)
			}
		})
	}
	t.Run("a default written in words", func(t *testing.T) {
		err := safely(t, "applyDefaults(&typo{})", func() error { return applyDefaults(&typo{}) })
		if err == nil {
			t.Errorf(`applyDefaults(&typo{}) error = nil, want setText's error for "twelve"`)
		}
	})
}
