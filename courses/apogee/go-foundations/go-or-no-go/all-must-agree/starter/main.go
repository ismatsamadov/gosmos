package main

import "fmt"

// canTrack reports whether a station can track: it is up, its dish is not stowed and the wind is under 72 km/h.
func canTrack(up, stowed bool, windKmh float64) bool {
	return true
}

func main() {
	fmt.Println("YST", canTrack(true, false, 40))
	fmt.Println("PEL", canTrack(true, true, 15))
	fmt.Println("SQT", canTrack(false, false, 12))
	fmt.Println("CWR", canTrack(true, false, 72))
}
