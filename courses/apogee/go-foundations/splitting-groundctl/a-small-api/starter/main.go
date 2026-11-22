package main

import (
	"fmt"

	"apogee/groundctl/telemetry"
)

func main() {
	f, err := telemetry.Parse("thistle 4182 2027-05-11T04:13:07Z batt_v=7.91 temp_c=-4.5 mode=NOMINAL")
	if err != nil {
		fmt.Println("rejected:", err)
		return
	}
	fmt.Printf("%s %d: battery %.2f V\n", f.SatID, f.Seq, f.BattV)
}
