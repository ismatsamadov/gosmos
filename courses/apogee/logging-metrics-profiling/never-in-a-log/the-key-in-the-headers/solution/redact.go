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
