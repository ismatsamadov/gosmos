package main

import (
	"errors"
	"fmt"
)

// minBattV is the lowest battery voltage ops accepts.
const minBattV = 6.6

// checkBattery returns an error when the battery reading v is under minBattV.
func checkBattery(v float64) error {
	if v < minBattV {
		return errors.New("battery too low")
	}
	return nil
}

func main() {
	readings := []float64{7.91, 6.20, 7.40}
	stored := 0
	for _, v := range readings {
		checkBattery(v)
		fmt.Printf("%.2f: stored\n", v)
		stored++
	}
	fmt.Printf("stored %d of %d\n", stored, len(readings))
}
