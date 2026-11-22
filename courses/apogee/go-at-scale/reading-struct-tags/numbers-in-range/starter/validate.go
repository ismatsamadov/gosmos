package main

import (
	"fmt"
	"reflect"
	"strings"
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
	return 0, nil, nil, false
}
