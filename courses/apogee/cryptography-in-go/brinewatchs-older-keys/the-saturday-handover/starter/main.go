package main

import "crypto/rsa"

// handOver checks Brinewatch's request, then returns the two lines that carry its new secrets,
// each locked to the receiving key the request names.
func handOver(signKey *rsa.PublicKey, r Request, apiKey, webhook []byte) ([]string, error) {
	return nil, nil
}
