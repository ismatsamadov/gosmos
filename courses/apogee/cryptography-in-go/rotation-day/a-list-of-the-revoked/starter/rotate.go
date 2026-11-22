package main

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"math/big"
	"time"
)

// revoke signs revocation list number, naming entries, with ca's key: in force from now for 7 days.
func revoke(ca Issuer, number int64, entries []x509.RevocationListEntry, now time.Time) ([]byte, error) {
	_, _, _ = rand.Reader, big.NewInt, fmt.Errorf
	return nil, nil
}

// isRevoked reports whether cert is on the list crlDER, which must be signed by ca and current at now.
func isRevoked(crlDER []byte, ca, cert *x509.Certificate, now time.Time) (bool, error) {
	return false, nil
}
