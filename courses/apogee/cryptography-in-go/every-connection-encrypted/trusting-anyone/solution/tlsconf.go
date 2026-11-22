package main

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
)

// newClient returns the client uplinkd reaches the station host with, trusting roots and judging dates by now.
func newClient(roots *x509.CertPool, host string, now func() time.Time) *http.Client {
	return &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    roots,
				ServerName: host,
				Time:       now,
			},
		},
	}
}
