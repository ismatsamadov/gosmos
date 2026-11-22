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
	v, _ := strconv.ParseFloat(text, 64)
	return v, nil
}

func main() {
	v, err := battery("thistle 4183 2027-05-11T04:13:37Z batt_v=7.9l temp_c=-4.5 mode=NOMINAL")
	fmt.Println(v, err)
}
