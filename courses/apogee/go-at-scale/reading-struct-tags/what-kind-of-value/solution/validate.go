package main

import (
	"fmt"
	"reflect"
)

// fieldLines returns a line for each field of the struct v: its
// name, its type and its kind.
func fieldLines(v any) []string {
	var lines []string
	for f := range reflect.TypeOf(v).Fields() {
		kind := f.Type.Kind()
		line := fmt.Sprintf("%s %s %s", f.Name, f.Type, kind)
		lines = append(lines, line)
	}
	return lines
}
