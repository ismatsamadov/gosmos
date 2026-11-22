package main

import (
	"fmt"
	"reflect"
)

// checkRule checks v, the value of the field called name, against
// one rule of its check tag.
func checkRule(name string, v reflect.Value, rule string) error {
	return fmt.Errorf("%s: unknown rule %q", name, rule)
}
