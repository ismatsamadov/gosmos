package main

import (
	"crypto/ed25519"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
)

// uplinkClient returns the client uplinkd queues commands with at the station host, proving itself with cert.
func uplinkClient(roots *x509.CertPool, host string, cert tls.Certificate, now func() time.Time) *http.Client {
	return &http.Client{}
}

// stationServer returns the config of a station presenting leaf, issued by inter, with key, to clients that chain to roots.
func stationServer(leaf, inter *x509.Certificate, key ed25519.PrivateKey, roots *x509.CertPool, now func() time.Time) *tls.Config {
	return &tls.Config{}
}
