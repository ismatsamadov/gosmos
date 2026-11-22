package main

import "fmt"

// batteryBand sorts a battery voltage into FULL, OK, LOW or CRITICAL.
func batteryBand(v float64) string {
	switch {
	case v >= 8.2:
		return "FULL"
	case v >= 7.2:
		return "OK"
	case v >= 6.6:
		return "LOW"
	default:
		return "CRITICAL"
	}
}

func main() {
	fmt.Println(batteryBand(7.91))
}
