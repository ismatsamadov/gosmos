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
	block, _ := pem.Decode(data)
	if block == nil {
		fmt.Println("yst.pem holds no PEM block")
		return
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("subject:", cert.Subject.CommonName)
	fmt.Println("issuer:", cert.Issuer.CommonName)
	fmt.Println("serial:", cert.SerialNumber)
	fmt.Println("valid:", cert.NotBefore.Format(time.DateOnly), "to", cert.NotAfter.Format(time.DateOnly))
	fmt.Println("names:", cert.DNSNames)
	fmt.Println("CA:", cert.IsCA)
}
