package main

import "fmt"

// aosMinute returns the first minute Pipit-3 is at or above maskDeg, or -1 if it never is.
func aosMinute(maskDeg float64) int {
	aos := -1
	for minute := range 15 {
		if elevationAt(minute) >= maskDeg {
			aos = minute
			break
		}
	}
	return aos
}

func main() {
	fmt.Println("AOS over Ystmark at minute", aosMinute(5))
}
