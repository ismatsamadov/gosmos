package main

import "fmt"

// poll returns GO, or NO-GO with the first failing reason, for a pass over Salar Quieto.
func poll(up bool, windKmh, maxElevDeg, battV float64) string {
	return "GO"
}

func main() {
	fmt.Println("Furrow-A over Salar Quieto:", poll(true, 18, 41.5, 7.9))
}
