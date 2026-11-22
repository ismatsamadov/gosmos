package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/big"
	"slices"
	"time"
)

// signingKey returns the version uplinkd seals with at now: the newest by Activate that has activated and is neither retired nor revoked.
func signingKey(ring []KeyVersion, now time.Time) (KeyVersion, error) {
	var best KeyVersion
	found := false
	for _, v := range ring {
		usable := !v.Revoked && !now.Before(v.Activate) && (v.Retire.IsZero() || now.Before(v.Retire))
		if usable && (!found || v.Activate.After(best.Activate)) {
			best, found = v, true
		}
	}
	if !found {
		return KeyVersion{}, fmt.Errorf("no key to sign with at %s", now.UTC().Format(time.RFC3339))
	}
	return best, nil
}

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

// withRevocations makes cfg refuse a peer whose certificate is on the list current returns, judged at now.
func withRevocations(cfg *tls.Config, current func() []byte, ca *x509.Certificate, now func() time.Time) {
	cfg.VerifyConnection = func(cs tls.ConnectionState) error {
		list := current()
		if list == nil {
			return nil // no list published yet
		}
		leaf := cs.PeerCertificates[0]
		revoked, err := isRevoked(list, ca, leaf, now())
		if err != nil {
			return err
		}
		if revoked {
			return fmt.Errorf("certificate %v is revoked", leaf.SerialNumber)
		}
		return nil
	}
}

// switchTime returns the earliest start, from from on, of a 10-minute certificate switch at station that keeps 5 minutes clear of every pass there.
func switchTime(passes []Pass, station string, from time.Time) (time.Time, error) {
	var mine []Pass
	for _, p := range passes {
		if p.Station == station {
			mine = append(mine, p)
		}
	}
	slices.SortFunc(mine, func(a, b Pass) int { return a.AOS.Compare(b.AOS) })
	start := from
	for _, p := range mine {
		if !start.Add(10 * time.Minute).After(p.AOS.Add(-5 * time.Minute)) {
			break // the switch ends before this pass's clear zone starts
		}
		if clear := p.LOS.Add(5 * time.Minute); start.Before(clear) {
			start = clear
		}
	}
	if end := from.Add(24 * time.Hour); start.Add(10 * time.Minute).After(end) {
		return time.Time{}, fmt.Errorf("no 10-minute gap for %s before %s", station, end.Format(time.RFC3339))
	}
	return start, nil
}
