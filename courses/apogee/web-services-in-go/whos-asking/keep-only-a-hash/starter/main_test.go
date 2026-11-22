package main

import "testing"

// The beta customers' test API keys, as Yusra gave them out. The service keeps none of them.
const (
	apogeeKey      = "apogee.9e4d2a7c1b0f3e5d8a6c4b2f0e1d3a5c"
	furrowlightKey = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"
	brinewatchKey  = "brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0"
	cirrineKey     = "cirrine.c3a1f7e95d2b4086a7e1c9f3b5d8a024"
)

// customers is the customer table as the service keeps it: a hash of each key, never the key.
var customers = map[string]Customer{
	"apogee":      {ID: "apogee", Name: "Apogee", KeyHash: "8018b533030c734641dda97af10c84dd57357d0e37b7a0bdee7ed2a63c911407"},
	"furrowlight": {ID: "furrowlight", Name: "Furrowlight", KeyHash: "95ff036b9ad2e59c26fae1be7dda78b12412d243fcbb6567bc8124e6f1c93fd3"},
	"brinewatch":  {ID: "brinewatch", Name: "Brinewatch", KeyHash: "6f5748f2886bd91e86930360f00a38ae260f6de18ea20012f4e4ec718d483f85"},
	"cirrine":     {ID: "cirrine", Name: "Cirrine Atmospherics", KeyHash: "d9ee2e5de1796bbeba223c0737933b3e60d53687d479f8a04c98c99df87b510d"},
}

func TestCheckKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string // the customer's ID, or "" when the key must be refused
	}{
		{"Furrowlight's key", furrowlightKey, "furrowlight"},
		{"Apogee's own key", apogeeKey, "apogee"},
		{"Brinewatch's key", brinewatchKey, "brinewatch"},
		{"Cirrine's key", cirrineKey, "cirrine"},
		{"a wrong secret", "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6c", ""},
		{"a customer the service doesn't know", "kessick.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b", ""},
		{"the secret in capitals", "furrowlight.5B0E7C1D9A3F4E6B8C2D1A0F9E8D7C6B", ""},
		{"Furrowlight's hash, stolen and sent as a key", customers["furrowlight"].KeyHash, ""},
		{"an empty key", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := checkKey(customers, tt.key)
			if got != customers[tt.want] || ok != (tt.want != "") {
				t.Errorf("checkKey(customers, %q) = %+v, %v, want %+v, %v", tt.key, got, ok, customers[tt.want], tt.want != "")
			}
		})
	}
}
