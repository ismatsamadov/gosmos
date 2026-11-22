package main

import (
	"crypto/hkdf"
	"crypto/sha256"
	"fmt"
)

func main() {
	master := []byte("uplink test master key, 32 bytes") // the tests' k-2028-07
	salt := []byte("apogee uplink keys v1")
	labels := []string{"sign|pipit-2", "sign|furrow-a", "sign|cirrine-1", "seal|cirrine-1"}
	for _, label := range labels {
		key, err := hkdf.Key(sha256.New, master, salt, "uplink", 32)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%-14s %x\n", label, key[:4])
	}
}
