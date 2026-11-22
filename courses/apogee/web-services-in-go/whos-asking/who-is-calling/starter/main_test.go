package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// furrowlightKey is Furrowlight's test API key.
const furrowlightKey = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"

// furrowlightToken is a token of Furrowlight's, as the service will hand them out.
const furrowlightToken = "furrowlight.1830682800.767de44da59ef9b9e3afb9354270dd7cc56e26ceb00650b8a615afe1cd10b616"

// describe writes what credentials returned: the credential, then nil or the APIError's status, code and message.
func describe(credential string, err error) string {
	if err == nil {
		return fmt.Sprintf("%q, nil", credential)
	}
	if e, ok := errors.AsType[*APIError](err); ok {
		return fmt.Sprintf("%q, %d %s %q", credential, e.Status, e.Code, e.Message)
	}
	return fmt.Sprintf("%q, %q (not an *APIError)", credential, err)
}

func TestCredentials(t *testing.T) {
	tests := []struct {
		name   string
		header string // "" sends no Authorization header
		scheme string
		want   string // "" when the request must be refused
	}{
		{"Furrowlight's key", "Key " + furrowlightKey, "Key", furrowlightKey},
		{"the scheme in lowercase", "key " + furrowlightKey, "Key", furrowlightKey},
		{"the scheme in capitals", "KEY " + furrowlightKey, "Key", furrowlightKey},
		{"a token where a token belongs", "Bearer " + furrowlightToken, "Bearer", furrowlightToken},
		{"no header at all", "", "Key", ""},
		{"a token where a key belongs", "Bearer x", "Key", ""},
		{"the scheme alone", "Key", "Key", ""},
		{"nothing after the space", "Key ", "Key", ""},
		{"no space after the scheme", "Key" + furrowlightKey, "Key", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/tokens", nil)
			call := fmt.Sprintf("credentials(no Authorization header, %q)", tt.scheme)
			if tt.header != "" {
				r.Header.Set("Authorization", tt.header)
				call = fmt.Sprintf("credentials(Authorization: %q, %q)", tt.header, tt.scheme)
			}
			want := describe(tt.want, nil)
			if tt.want == "" {
				want = `"", 401 unauthorized "missing credentials"`
			}
			if got := describe(credentials(r, tt.scheme)); got != want {
				t.Errorf("%s = %s, want %s", call, got, want)
			}
		})
	}
}
