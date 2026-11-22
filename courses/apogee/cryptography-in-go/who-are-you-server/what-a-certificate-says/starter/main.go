package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

func main() {
	data, err := os.ReadFile("yst.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _, _ = pem.Decode, x509.ParseCertificate, time.DateOnly
	fmt.Printf("yst.pem: %d bytes\n", len(data))
}
