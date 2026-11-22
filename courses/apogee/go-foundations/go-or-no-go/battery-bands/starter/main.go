package main

import "fmt"

// batteryLine returns a battery's ops band, FULL, OK, LOW or CRITICAL, followed by its voltage.
func batteryLine(battV float64) string {
	var band string
	// TODO: set band from battV
	return fmt.Sprintf("%s (%.2f V)", band, battV)
}

func main() {
	fmt.Println("Pipit-1 battery:", batteryLine(7.91))
}
