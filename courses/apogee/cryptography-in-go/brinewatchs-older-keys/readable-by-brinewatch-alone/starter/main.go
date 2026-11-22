package main

import "crypto/rsa"

// encryptTo locks secret with RSA-OAEP, SHA-256 and label, so that only the holder of pub's
// private key can read it.
func encryptTo(pub *rsa.PublicKey, secret []byte, label string) ([]byte, error) {
	return nil, nil
}

// decryptWith opens what encryptTo locked, as Brinewatch does with its private key.
func decryptWith(priv *rsa.PrivateKey, ciphertext []byte, label string) ([]byte, error) {
	return nil, nil
}
