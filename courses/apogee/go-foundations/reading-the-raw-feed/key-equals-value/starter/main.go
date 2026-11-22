package main

import (
	"fmt"
	"strings"
)

// field returns the value of key in a telemetry line, and whether the line has that key.
func field(line, key string) (string, bool) {
	for _, f := range strings.Fields(line) {
		_ = f // TODO: cut f at its first "=" and return the value when the key matches.
	}
	return "", false
}

func main() {
	line := "thistle 4182 2027-05-11T04:13:07Z batt_v=7.91 temp_c=-4.5 mode=NOMINAL"
	value, ok := field(line, "batt_v")
	fmt.Println(value, ok)
}
