package main

import (
	"fmt"
	"reflect"
)

// describe returns name=value for each exported field of the
// struct v, or of the struct v points to; anything else gives nil.
func describe(v any) []string {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}
	var lines []string
	for f, fv := range rv.Fields() {
		if f.IsExported() {
			lines = append(lines, fmt.Sprintf("%s=%v", f.Name, fv))
		}
	}
	return lines
}
