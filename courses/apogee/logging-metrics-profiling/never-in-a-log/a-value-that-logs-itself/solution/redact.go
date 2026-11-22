package main

import (
	"log/slog"
	"net/http"
	"strings"
)

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

// LogValue logs the key as its customer ID and [redacted], never its secret.
func (k APIKey) LogValue() slog.Value {
	customer, _, ok := strings.Cut(string(k), ".")
	if !ok {
		return slog.StringValue("[redacted]")
	}
	return slog.StringValue(customer + ".[redacted]")
}
