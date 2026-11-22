package main

import (
	"fmt"
	"math"
)

// bits returns the strength, in bits, of a secret of length symbols from an alphabet of the given size.
func bits(length, alphabet int) float64 {
	return float64(length) * math.Log2(float64(alphabet))
}

func main() {
	fmt.Printf("6 digits: %.1f bits\n", bits(6, 10))
	fmt.Printf("4 bytes in hex: %.1f bits\n", bits(8, 16)) // 8 digits, 0-9 and a-f
	fmt.Printf("rand.Text: %.1f bits\n", bits(26, 32))     // 26 characters, A-Z and 2-7
	fmt.Printf("32-byte key: %.1f bits\n", bits(32, 256))  // 32 bytes of 256 values each
}
