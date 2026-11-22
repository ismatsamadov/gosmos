package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
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

// setText sets v from text, read as v's type reads it: a duration
// as duration text, an int or a float as a number, a string as it
// is written.
func setText(v reflect.Value, text string) error {
	switch {
	case v.Type() == reflect.TypeFor[time.Duration]():
		d, err := time.ParseDuration(text)
		if err != nil {
			return err
		}
		v.SetInt(int64(d))
	case v.CanInt():
		n, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case v.CanFloat():
		x, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return err
		}
		v.SetFloat(x)
	case v.Kind() == reflect.String:
		v.SetString(text)
	default:
		return fmt.Errorf("no default for %s", v.Type())
	}
	return nil
}
