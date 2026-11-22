package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

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
