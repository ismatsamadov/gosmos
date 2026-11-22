package main

import (
	"fmt"
	"reflect"
)

// tagLines returns a line for each field of the struct v: its
// name and its check rules.
func tagLines(v any) []string {
	var lines []string
	for f := range reflect.TypeOf(v).Fields() {
		check, ok := f.Tag.Lookup("check")
		if !ok {
			check = "no check"
		}
		lines = append(lines, fmt.Sprintf("%s: %s", f.Name, check))
	}
	return lines
}
