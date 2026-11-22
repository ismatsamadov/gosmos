package main

import "fmt"

// windCall returns what Cape Wren's dish does in a wind of windKmh: track, track with warning, or stow.
func windCall(windKmh float64) string {
	return ""
}

func main() {
	fmt.Println("Cape Wren at 64 km/h:", windCall(64))
}
