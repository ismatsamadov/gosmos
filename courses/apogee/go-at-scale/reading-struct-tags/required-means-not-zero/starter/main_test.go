package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

// icmFirst is Ilha Calmaria's form as its host first sent it.
var icmFirst = HostConfig{Code: "ICM", LatDeg: 8, LonDeg: -30, MaskDeg: 7, Dishes: 1, Bands: []string{"S", "K"}}

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
		{"Ilha Calmaria's first form, no name", icmFirst, "Name", "required", "HostConfig.Name: required"},
		{"the name written in", HostConfig{Name: "Ilha Calmaria"}, "Name", "required", ""},
		{"no bands at all", HostConfig{}, "Bands", "required", "HostConfig.Bands: required"},
		{"bands cleared to an empty list", HostConfig{Bands: []string{}}, "Bands", "required", "HostConfig.Bands: required"},
		{"one band, S", HostConfig{Bands: []string{"S"}}, "Bands", "required", ""},
		{"no dishes", HostConfig{}, "Dishes", "required", "HostConfig.Dishes: required"},
		{"two dishes", HostConfig{Dishes: 2}, "Dishes", "required", ""},
		{"a rule nobody knows", HostConfig{Name: "Ilha Calmaria"}, "Name", "sometimes", `HostConfig.Name: unknown rule "sometimes"`},
	})
}
