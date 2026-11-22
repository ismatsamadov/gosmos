package main

import "fmt"

// averageBattery returns the mean reading from minute from up to, not including, minute to, skipping lost frames.
func averageBattery(from, to int) float64 {
	sum, readings := 0.0, 0
	for minute := from; minute < to; minute++ {
		battV := batteryAt(minute)
		if battV == lostFrame {
			continue
		}
		sum += battV
		readings++
	}
	if readings == 0 {
		return 0
	}
	return sum / float64(readings)
}

func main() {
	fmt.Printf("Average battery: %.2f V\n", averageBattery(0, 6))
}
