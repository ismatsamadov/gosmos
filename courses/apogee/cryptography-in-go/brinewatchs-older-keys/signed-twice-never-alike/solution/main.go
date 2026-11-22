package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// signRequest signs req with RSA-PSS over its SHA-256 digest, as Brinewatch's security office
// signs its requests.
func signRequest(priv *rsa.PrivateKey, req []byte) ([]byte, error) {
	digest := sha256.Sum256(req)
	return rsa.SignPSS(rand.Reader, priv, crypto.SHA256, digest[:], nil)
}

// checkRequest returns nil only if sig is pub's RSA-PSS signature over req's SHA-256 digest.
func checkRequest(pub *rsa.PublicKey, req, sig []byte) error {
	digest := sha256.Sum256(req)
	return rsa.VerifyPSS(pub, crypto.SHA256, digest[:], sig, nil)
}
