package main

import "fmt"

// batteryLine returns a battery's ops band, FULL, OK, LOW or CRITICAL, followed by its voltage.
func batteryLine(battV float64) string {
	var band string
	if battV >= 8.2 {
		band = "FULL"
	} else if battV >= 7.2 {
		band = "OK"
	} else if battV >= 6.6 {
		band = "LOW"
	} else {
		band = "CRITICAL"
	}
	return fmt.Sprintf("%s (%.2f V)", band, battV)
}

func main() {
	fmt.Println("Pipit-1 battery:", batteryLine(7.91))
}
