package main

import "fmt"

// poll returns GO, or NO-GO with the first failing reason, for a pass over Salar Quieto.
func poll(up bool, windKmh, maxElevDeg, battV float64) string {
	switch {
	case !up:
		return "NO-GO: station down"
	case windKmh >= 72:
		return "NO-GO: wind"
	case maxElevDeg < 3:
		return "NO-GO: below mask"
	case battV < 6.6:
		return "NO-GO: battery"
	default:
		return "GO"
	}
}

func main() {
	fmt.Println("Furrow-A over Salar Quieto:", poll(true, 18, 41.5, 7.9))
}
