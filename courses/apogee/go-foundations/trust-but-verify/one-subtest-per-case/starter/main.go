package main

import "fmt"

// windCall returns what Cape Wren's dish does in a wind of windKmh: track, track with warning, or stow.
func windCall(windKmh float64) string {
	switch {
	case windKmh >= 72:
		return "stow"
	case windKmh >= 50:
		return "track with warning"
	default:
		return "track"
	}
}

func main() {
	fmt.Println(windCall(64))
}
