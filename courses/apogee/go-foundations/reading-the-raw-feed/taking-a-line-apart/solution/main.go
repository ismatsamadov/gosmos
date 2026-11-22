package main

import (
	"fmt"
	"strings"
)

// line is a telemetry line from Thistle, exactly as it came down.
const line = "thistle  4182 2027-05-11T04:13:07Z batt_v=7.91 temp_c=-4.5 mode=NOMINAL"

func main() {
	fields := strings.Fields(line)
	fmt.Println("fields:", len(fields))
	fmt.Println("satellite:", fields[0])
	fmt.Println("seq:", fields[1])
}
