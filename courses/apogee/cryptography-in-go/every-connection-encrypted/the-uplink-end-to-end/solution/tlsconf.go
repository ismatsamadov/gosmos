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
	return &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{
			RootCAs:      roots,
			ServerName:   host,
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
			Time:         now,
		}},
	}
}

// stationServer returns the config of a station presenting leaf, issued by inter, with key, to clients that chain to roots.
func stationServer(leaf, inter *x509.Certificate, key ed25519.PrivateKey, roots *x509.CertPool, now func() time.Time) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{leaf.Raw, inter.Raw},
			PrivateKey:  key,
			Leaf:        leaf,
		}},
		MinVersion: tls.VersionTLS13,
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  roots,
		Time:       now,
	}
}
