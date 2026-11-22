package main

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

// unauthorized returns the 401 that tells the caller msg.
func unauthorized(msg string) error {
	return &APIError{Status: http.StatusUnauthorized, Code: "unauthorized", Message: msg}
}

// credentials returns the credential that follows scheme in r's Authorization header.
func credentials(r *http.Request, scheme string) (string, error) {
	got, credential, _ := strings.Cut(r.Header.Get("Authorization"), " ")
	if !strings.EqualFold(got, scheme) || credential == "" {
		return "", unauthorized("missing credentials")
	}
	return credential, nil
}

// checkKey returns the ID of the customer whose key this is, if keys holds exactly this key for it.
func checkKey(keys map[string]string, key string) (string, bool) {
	id, _, _ := strings.Cut(key, ".")
	want, ok := keys[id]
	if !ok || subtle.ConstantTimeCompare([]byte(key), []byte(want)) != 1 {
		return "", false
	}
	return id, true
}
