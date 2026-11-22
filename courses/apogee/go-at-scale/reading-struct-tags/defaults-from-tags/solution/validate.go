package main

import (
	"fmt"
	"reflect"
)

// applyDefaults sets each zero exported field of the struct ptr
// points to from the field's default tag.
func applyDefaults(ptr any) error {
	v := reflect.Indirect(reflect.ValueOf(ptr))
	if v.Kind() != reflect.Struct || !v.CanSet() {
		return fmt.Errorf(
			"defaults: want a pointer to a struct, got %T", ptr)
	}
	for f, fv := range v.Fields() {
		text, ok := f.Tag.Lookup("default")
		if !ok || !fv.CanSet() || !fv.IsZero() {
			continue
		}
		if err := setText(fv, text); err != nil {
			return fmt.Errorf("defaults: %s: %w", f.Name, err)
		}
	}
	return nil
}
