package main

import (
	"errors"
	"fmt"
)

// minBattV is the lowest battery voltage ops accepts.
const minBattV = 6.6

// checkBattery returns an error when satellite sat's battery reading v is under minBattV.
func checkBattery(sat string, v float64) error {
	if v < minBattV {
		return errors.New("battery too low")
	}
	return nil
}

func main() {
	fmt.Println(checkBattery("thistle", 6.2))
}
