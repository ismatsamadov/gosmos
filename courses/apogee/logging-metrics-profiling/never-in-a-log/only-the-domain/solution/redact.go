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

// LogValue logs the key as its String method prints it, so slog's handlers never write its secret.
func (k Key) LogValue() slog.Value {
	return slog.StringValue(k.String())
}

// redact replaces the value of every attribute whose key names a secret with [redacted], keeping the key.
func redact(groups []string, a slog.Attr) slog.Attr {
	switch strings.ToLower(a.Key) {
	case "api_key", "authorization", "password", "secret", "token":
		return slog.String(a.Key, "[redacted]")
	}
	return a
}

// maskEmail keeps an address's first character and its domain, as in h***@furrowlight.example.
func maskEmail(s string) string {
	local, domain, ok := strings.Cut(s, "@")
	if !ok || local == "" || domain == "" || strings.Contains(domain, "@") {
		return "[redacted]"
	}
	return string([]rune(local)[:1]) + "***@" + domain
}
