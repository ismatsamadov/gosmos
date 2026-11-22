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
	fmt.Println(checkBattery(6.2))
}
