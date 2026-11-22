package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	texts := []struct{ name, text string }{
		{"Pipit-1", "Apogee ground calling Pipit-1"},
		{"Pipit-2", "Apogee ground calling Pipit-2"},
	}
	fmt.Printf("digest size: %d bytes\n", sha256.Size)
	for _, t := range texts {
		sum := sha256.Sum256([]byte(t.text))
		fmt.Printf("%s: %x\n", t.name, sum[:8])
	}
}
