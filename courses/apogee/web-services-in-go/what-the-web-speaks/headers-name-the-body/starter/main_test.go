package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIsJSON(t *testing.T) {
	tests := []struct {
		name   string
		header string // the name the header is set under; "" leaves it out
		value  string
		want   bool
	}{
		{"the stand-in's answer", "Content-Type", "application/json", true},
		{"Cirrine's good answer", "Content-Type", "application/json; charset=utf-8", true},
		{"JSON in capitals", "Content-Type", "Application/JSON", true},
		{"a name set in lowercase", "content-type", "application/json", true},
		{"a plain-text error", "Content-Type", "text/plain; charset=utf-8", false},
		{"JSONP, which only starts the same", "Content-Type", "application/jsonp", false},
		{"no Content-Type at all", "", "", false},
		{"a value cut short", "Content-Type", "application/", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			call := "isJSON(no Content-Type)"
			if tt.header != "" {
				h.Set(tt.header, tt.value)
				call = fmt.Sprintf("isJSON(%s: %q)", tt.header, tt.value)
			}
			if got := isJSON(h); got != tt.want {
				t.Errorf("%s = %v, want %v", call, got, tt.want)
			}
		})
	}
}
