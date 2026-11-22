package main

import "net/http"

// safeHeader returns a copy of h in which every value of a secret header reads [redacted].
func safeHeader(h http.Header) http.Header {
	c := h.Clone()
	for _, name := range []string{"Authorization", "Cookie", "Set-Cookie", "Apogee-Signature"} {
		for i := range c[name] {
			c[name][i] = "[redacted]"
		}
	}
	return c
}

// APIKey is a customer's API key: its customer ID, a dot, then 32 hex digits of secret.
type APIKey string

// TODO: give APIKey a LogValue method that logs its customer ID and [redacted], never its secret
