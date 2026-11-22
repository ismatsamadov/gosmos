package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

// vrlForm is Vorrel Island's form as its host sent it.
var vrlForm = HostConfig{Code: "VRL", Name: "Vorrel Island", LatDeg: -64, LonDeg: -140, MaskDeg: 4, Dishes: 0, Bands: []string{"S", "X"}}

// icmFixed is Ilha Calmaria's form as its host sent it again, its name written in.
var icmFixed = HostConfig{Code: "ICM", Name: "Ilha Calmaria", LatDeg: 8, LonDeg: -30, MaskDeg: 7, Dishes: 1, Bands: []string{"S", "K"}}

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
	firstWave := "oneof=TIR VSK CAL BRM OST"
	run(t, []ruleCase{
		{"Ilha Calmaria's bands, S and K", icmFixed, "Bands", "oneof=UHF S X", `HostConfig.Bands: "K" is not one of UHF S X`},
		{"Vorrel Island's bands, S and X", vrlForm, "Bands", "oneof=UHF S X", ""},
		{"Port Ellery's bands, UHF and S", HostConfig{Bands: []string{"UHF", "S"}}, "Bands", "oneof=UHF S X", ""},
		{"two bands typed in one box", HostConfig{Bands: []string{"S X"}}, "Bands", "oneof=UHF S X", `HostConfig.Bands: "S X" is not one of UHF S X`},
		{"a first-wave code", HostConfig{Code: "VSK"}, "Code", firstWave, ""},
		{"a third-wave code", vrlForm, "Code", firstWave, `HostConfig.Code: "VRL" is not one of TIR VSK CAL BRM OST`},
		{"Vorrel Island's code", vrlForm, "Code", "len=3", ""},
		{"a code one letter too long", HostConfig{Code: "VRLX"}, "Code", "len=3", `HostConfig.Code: "VRLX" is not 3 characters`},
		{"a tag that asks for four letters", vrlForm, "Code", "len=4", `HostConfig.Code: "VRL" is not 4 characters`},
		{"a code with an accent, 3 characters in 4 bytes", HostConfig{Code: "ÍCM"}, "Code", "len=3", ""},
		{"tlmd's feeds, one cut short", PipelineConfig{Feeds: []string{"TIR", "YS"}}, "Feeds", "len=3", `PipelineConfig.Feeds: "YS" is not 3 characters`},
		{"the first wave's feeds", PipelineConfig{Feeds: []string{"TIR", "VSK", "CAL", "BRM", "OST"}}, "Feeds", "len=3", ""},
		{"a length on the slot count", Config{Slots: 12}, "Slots", "len=3", "Config.Slots: rule len=3 does not fit int"},
		{"0 dishes, still below min", vrlForm, "Dishes", "min=1", "HostConfig.Dishes: 0 is below min 1"},
	})
}
