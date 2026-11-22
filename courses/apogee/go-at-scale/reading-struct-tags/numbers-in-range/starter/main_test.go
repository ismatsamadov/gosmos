package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

// vrlForm is Vorrel Island's form as its host sent it.
var vrlForm = HostConfig{Code: "VRL", Name: "Vorrel Island", LatDeg: -64, LonDeg: -140, MaskDeg: 4, Dishes: 0, Bands: []string{"S", "X"}}

// friday is the Config the test deploy of Friday 1 June read.
var friday = Config{Slots: 0, Days: 100, MinElevDeg: 10, Turnaround: 2 * time.Minute, Replan: 10 * time.Minute, RunLimit: 30 * time.Second, Queue: 256, Quiet: 5 * time.Second}

// archiveLimits are the pass archive's settings: its size is an int64
// of bytes, not a duration.
type archiveLimits struct {
	MaxBytes int64
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

// show writes a field's value as the messages here quote it.
func show(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Slice:
		if v.IsNil() {
			return "nil"
		}
		return fmt.Sprintf("%q", v.Interface())
	}
	return fmt.Sprint(v.Interface())
}

// message is err's text in quotes, or nil.
func message(err error) string {
	if err == nil {
		return "nil"
	}
	return strconv.Quote(err.Error())
}

// quoted is a problem in quotes, or nil for none.
func quoted(problem string) string {
	if problem == "" {
		return "nil"
	}
	return strconv.Quote(problem)
}

// ruleCase is one call of checkRule on a field of a form: want is the
// problem it reports, or "" for none.
type ruleCase struct {
	name  string
	form  any
	field string
	rule  string
	want  string
}

// run checks each case as a subtest.
func run(t *testing.T, cases []ruleCase) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := reflect.ValueOf(c.form).FieldByName(c.field)
			name := reflect.TypeOf(c.form).Name() + "." + c.field
			call := fmt.Sprintf("checkRule(%q, %s, %q)", name, show(v), c.rule)
			got := safely(t, call, func() error { return checkRule(name, v, c.rule) })
			if message(got) != quoted(c.want) {
				t.Errorf("%s = %s, want %s", call, message(got), quoted(c.want))
			}
		})
	}
}

func TestCheckRule(t *testing.T) {
	run(t, []ruleCase{
		{"Vorrel Island's form, 0 dishes", vrlForm, "Dishes", "min=1", "HostConfig.Dishes: 0 is below min 1"},
		{"the two dishes its host meant", HostConfig{Dishes: 2}, "Dishes", "max=2", ""},
		{"Friday's test deploy, 0 slots", friday, "Slots", "min=1", "Config.Slots: 0 is below min 1"},
		{"seventy slots", Config{Slots: 70}, "Slots", "max=64", "Config.Slots: 70 is above max 64"},
		{"sixty-four slots, the most allowed", Config{Slots: 64}, "Slots", "max=64", ""},
		{"a hundred and twenty-one days", Config{Days: 121}, "Days", "max=120", "Config.Days: 121 is above max 120"},
		{"a mask of 15.5°", HostConfig{MaskDeg: 15.5}, "MaskDeg", "max=15", "HostConfig.MaskDeg: 15.5 is above max 15"},
		{"Vorrel Island at 64° south", vrlForm, "LatDeg", "min=-90", ""},
		{"a 15-minute turnaround", Config{Turnaround: 15 * time.Minute}, "Turnaround", "max=10m", "Config.Turnaround: 15m0s is above max 10m0s"},
		{"no turnaround at all", Config{}, "Turnaround", "min=0s", ""},
		{"a run limit of half a second", Config{RunLimit: 500 * time.Millisecond}, "RunLimit", "min=1s", "Config.RunLimit: 500ms is below min 1s"},
		{"an empty archive, an int64 that isn't a duration", archiveLimits{}, "MaxBytes", "min=1", "archiveLimits.MaxBytes: 0 is below min 1"},
		{"a minimum on a station code", vrlForm, "Code", "min=1", "HostConfig.Code: rule min=1 does not fit string"},
		{"a slot bound mistyped as 6A", Config{Slots: 12}, "Slots", "max=6A", "Config.Slots: rule max=6A does not fit int"},
		{"a turnaround bound with no unit", Config{Turnaround: 2 * time.Minute}, "Turnaround", "max=10", "Config.Turnaround: rule max=10 does not fit time.Duration"},
		{"a mask bound in words", HostConfig{MaskDeg: 4}, "MaskDeg", "max=fifteen", "HostConfig.MaskDeg: rule max=fifteen does not fit float64"},
		{"a name still required", HostConfig{}, "Name", "required", "HostConfig.Name: required"},
	})
}
