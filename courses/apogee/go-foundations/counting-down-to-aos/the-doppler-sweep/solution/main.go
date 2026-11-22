package main

import "fmt"

// bestOffset returns the offset from fromHz to toHz, in steps of stepHz, with the strongest signal; the first wins a tie.
func bestOffset(fromHz, toHz, stepHz int) int {
	bestHz, bestDBm := fromHz, signalStrength(fromHz)
	for offsetHz := fromHz + stepHz; offsetHz <= toHz; offsetHz += stepHz {
		dBm := signalStrength(offsetHz)
		if dBm > bestDBm {
			bestHz, bestDBm = offsetHz, dBm
		}
	}
	return bestHz
}

func main() {
	fmt.Printf("Pipit-3 answers at %+d Hz\n", bestOffset(-12000, 12000, 500))
}
