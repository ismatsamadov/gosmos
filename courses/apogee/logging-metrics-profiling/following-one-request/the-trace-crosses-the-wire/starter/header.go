package main

import "strings"

// remoteParent returns the span an Apogee-Trace header names, "<trace id>-<span id>", with only
// its IDs set, or false when the header is missing or malformed.
func remoteParent(header string) (*Span, bool) {
	traceID, spanID, ok := strings.Cut(header, "-")
	if !ok || !lowerHex(traceID, 16) || !lowerHex(spanID, 8) {
		return nil, false
	}
	return &Span{TraceID: traceID, SpanID: spanID}, true
}

// lowerHex reports whether s is exactly n lowercase hex digits.
func lowerHex(s string, n int) bool {
	if len(s) != n {
		return false
	}
	for _, c := range s {
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return true
}
