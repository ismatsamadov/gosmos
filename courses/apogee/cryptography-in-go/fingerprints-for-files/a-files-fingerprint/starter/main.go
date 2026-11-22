package main

import (
	"crypto/sha256"
	"fmt"
	"hash/crc32"
)

func main() {
	_ = sha256.Sum256
	texts := []struct{ name, text string }{
		{"Pipit-1", "Apogee ground calling Pipit-1"},
		{"Pipit-2", "Apogee ground calling Pipit-2"},
	}
	fmt.Printf("digest size: %d bytes\n", crc32.Size)
	for _, t := range texts {
		fmt.Printf("%s: %08x\n", t.name, crc32.ChecksumIEEE([]byte(t.text)))
	}
}
