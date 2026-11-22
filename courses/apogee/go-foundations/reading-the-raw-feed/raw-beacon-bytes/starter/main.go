package main

import "fmt"

// decodeBeacon reads the sequence number and battery millivolts from a Pipit beacon in hex.
func decodeBeacon(hexText string) (seq, battMV uint16, err error) {
	return 0, 0, nil
}

func main() {
	seq, battMV, err := decodeBeacon("4150011056001ee6ffd30082")
	fmt.Println(seq, battMV, err)
}
