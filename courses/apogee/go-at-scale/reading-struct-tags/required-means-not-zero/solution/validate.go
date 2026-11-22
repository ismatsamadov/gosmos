package main

import (
	"fmt"
	"reflect"
)

// checkRule checks v, the value of the field called name, against
// one rule of its check tag.
func checkRule(name string, v reflect.Value, rule string) error {
	switch rule {
	case "required":
		if v.IsZero() || v.Kind() == reflect.Slice && v.Len() == 0 {
			return fmt.Errorf("%s: required", name)
		}
		return nil
	}
	return fmt.Errorf("%s: unknown rule %q", name, rule)
}
