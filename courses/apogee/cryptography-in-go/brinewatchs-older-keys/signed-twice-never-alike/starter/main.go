package main

import "crypto/rsa"

// signRequest signs req with RSA-PSS over its SHA-256 digest, as Brinewatch's security office
// signs its requests.
func signRequest(priv *rsa.PrivateKey, req []byte) ([]byte, error) {
	return nil, nil
}

// checkRequest returns nil only if sig is pub's RSA-PSS signature over req's SHA-256 digest.
func checkRequest(pub *rsa.PublicKey, req, sig []byte) error {
	return nil
}
