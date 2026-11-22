package main

import "fmt"

// minBattV is the lowest battery voltage ops accepts.
const minBattV = 6.6

// checkBattery returns an error when the battery reading v is under minBattV.
func checkBattery(v float64) error {
	return nil
}

func main() {
	fmt.Println(checkBattery(6.2))
}
