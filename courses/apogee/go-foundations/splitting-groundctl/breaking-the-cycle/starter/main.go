package main

import (
	"fmt"

	"apogee/groundctl/alert"
	"apogee/groundctl/telemetry"
)

func main() {
	f := telemetry.Frame{SatID: "thistle", Seq: 4182, BattV: 7.05}
	fmt.Println(alert.Line(f, telemetry.Grade(f)))
}
