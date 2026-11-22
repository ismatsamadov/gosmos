package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// field returns the value of key in a telemetry line, and whether the line has that key.
func field(line, key string) (string, bool) {
	for _, f := range strings.Fields(line) {
		k, v, found := strings.Cut(f, "=")
		if found && k == key {
			return v, true
		}
	}
	return "", false
}

// battery returns the battery voltage a telemetry line reports.
func battery(line string) (float64, error) {
	text, ok := field(line, "batt_v")
	if !ok {
		return 0, errors.New("no batt_v field")
	}
	v, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("batt_v: %w", err)
	}
	return v, nil
}

func main() {
	fmt.Println(battery("thistle 4182 2027-05-11T04:13:07Z batt_v=7.91 temp_c=-4.5 mode=NOMINAL"))
}
