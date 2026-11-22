package main

import "fmt"

// windCall says what Cape Wren's dish does in a wind of kmh km/h.
func windCall(kmh float64) string {
	switch {
	case kmh >= 72:
		return "stow"
	case kmh >= 50:
		return "track with warning"
	default:
		return "track"
	}
}

func main() {
	fmt.Println(windCall(64))
}
