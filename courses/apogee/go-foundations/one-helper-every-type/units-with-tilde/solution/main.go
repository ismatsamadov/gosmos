package main

import "fmt"

// Volts is a battery voltage.
type Volts float64

// Kilometers is an altitude or a distance.
type Kilometers float64

// KmPerHour is a wind speed, as Cape Wren's anemometer reports it.
type KmPerHour float32

// Real admits the number types Clamp works with.
type Real interface {
	~float32 | ~float64
}

// Clamp returns v, raised to lo or lowered to hi if it falls outside them.
func Clamp[T Real](v, lo, hi T) T {
	return min(max(v, lo), hi)
}

func main() {
	fmt.Println(Clamp(7.5, 6.6, 8.4))
}
