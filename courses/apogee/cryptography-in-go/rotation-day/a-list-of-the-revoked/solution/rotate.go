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
	tmpl := &x509.RevocationList{
		Number:                    big.NewInt(number),
		ThisUpdate:                now,
		NextUpdate:                now.AddDate(0, 0, 7),
		RevokedCertificateEntries: entries,
	}
	return x509.CreateRevocationList(rand.Reader, tmpl, ca.Cert, ca.Key)
}

// isRevoked reports whether cert is on the list crlDER, which must be signed by ca and current at now.
func isRevoked(crlDER []byte, ca, cert *x509.Certificate, now time.Time) (bool, error) {
	list, err := x509.ParseRevocationList(crlDER)
	if err != nil {
		return false, err
	}
	if err := list.CheckSignatureFrom(ca); err != nil {
		return false, err
	}
	if now.Before(list.ThisUpdate) || now.After(list.NextUpdate) {
		return false, fmt.Errorf("revocation list %v is not current at %s", list.Number, now.Format(time.RFC3339))
	}
	for _, e := range list.RevokedCertificateEntries {
		if e.SerialNumber.Cmp(cert.SerialNumber) == 0 {
			return true, nil
		}
	}
	return false, nil
}
