package main

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// checkRule checks v, the value of the field called name, against
// one rule of its check tag.
func checkRule(name string, v reflect.Value, rule string) error {
	key, arg, _ := strings.Cut(rule, "=")
	switch key {
	case "required":
		if v.IsZero() || v.Kind() == reflect.Slice && v.Len() == 0 {
			return fmt.Errorf("%s: required", name)
		}
		return nil
	case "min", "max":
		c, got, bound, ok := measure(v, arg)
		switch {
		case !ok:
			return misfit(name, rule, v.Type())
		case key == "min" && c < 0:
			return fmt.Errorf("%s: %v is below min %v",
				name, got, bound)
		case key == "max" && c > 0:
			return fmt.Errorf("%s: %v is above max %v",
				name, got, bound)
		}
		return nil
	case "oneof", "len":
		if v.Type() == reflect.TypeFor[[]string]() {
			for i := range v.Len() {
				err := checkRule(name, v.Index(i), rule)
				if err != nil {
					return err
				}
			}
			return nil
		}
		if v.Kind() != reflect.String {
			return misfit(name, rule, v.Type())
		}
		s, words := v.String(), strings.Fields(arg)
		if key == "oneof" && !slices.Contains(words, s) {
			return fmt.Errorf("%s: %q is not one of %s",
				name, s, arg)
		}
		n := strconv.Itoa(utf8.RuneCountInString(s))
		if key == "len" && n != arg {
			return fmt.Errorf("%s: %q is not %s characters",
				name, s, arg)
		}
		return nil
	}
	return fmt.Errorf("%s: unknown rule %q", name, rule)
}

// misfit is the problem of a rule on a field of type t, a type that
// has no use for the rule.
func misfit(name, rule string, t reflect.Type) error {
	return fmt.Errorf("%s: rule %s does not fit %s", name, rule, t)
}

// measure compares v with text, read as a bound of v's type. It
// returns -1, 0 or +1, then v and the bound as a message prints
// them, and false when the bound doesn't fit v.
func measure(v reflect.Value, text string) (int, any, any, bool) {
	switch {
	case v.Type() == reflect.TypeFor[time.Duration]():
		d, err := time.ParseDuration(text)
		got := time.Duration(v.Int())
		return cmp.Compare(got, d), got, d, err == nil
	case v.CanInt():
		n, err := strconv.ParseInt(text, 10, 64)
		return cmp.Compare(v.Int(), n), v.Int(), n, err == nil
	case v.CanFloat():
		x, err := strconv.ParseFloat(text, 64)
		return cmp.Compare(v.Float(), x), v.Float(), x, err == nil
	}
	return 0, nil, nil, false
}
