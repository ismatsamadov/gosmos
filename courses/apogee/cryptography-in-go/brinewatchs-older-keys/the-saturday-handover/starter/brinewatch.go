package main

import (
	"bytes"
	"crypto"
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
)

// The labels Brinewatch's two secrets are locked under, so a locked API key never opens as a
// webhook secret.
const (
	apiKeyLabel  = "apogee api key"
	webhookLabel = "apogee webhook secret"
)

// Request is Brinewatch's signed request for its new secrets, with the key it names.
type Request struct {
	Text []byte // "brinewatch handover sha256:<64 hex>": the exact bytes signed
	Sig  []byte // RSA-PSS over Text's SHA-256 digest, by Brinewatch's signing key
	Key  []byte // the receiving key Text names, as one PEM PUBLIC KEY block
}

// keyFingerprint returns pub's fingerprint as Apogee writes digests: "sha256:" and the SHA-256
// of the key's PKIX bytes, in hex.
func keyFingerprint(pub *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(der)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// readRSA reads the RSA public key in one PEM PUBLIC KEY block, refusing any other kind of key
// and any RSA key under 2048 bits.
func readRSA(pemText []byte) (*rsa.PublicKey, error) {
	block, rest := pem.Decode(pemText)
	if block == nil || block.Type != "PUBLIC KEY" || len(bytes.TrimSpace(rest)) != 0 {
		return nil, errors.New("want one PUBLIC KEY block")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch k := key.(type) {
	case *rsa.PublicKey:
		if k.N.BitLen() < 2048 {
			return nil, fmt.Errorf("want at least 2048 bits, got %d", k.N.BitLen())
		}
		return k, nil
	case ed25519.PublicKey:
		return nil, errors.New("want an RSA key, got Ed25519")
	case *ecdh.PublicKey:
		return nil, errors.New("want an RSA key, got X25519")
	default:
		return nil, fmt.Errorf("want an RSA key, got %T", k)
	}
}

// encryptTo locks secret with RSA-OAEP, SHA-256 and label, so that only the holder of pub's
// private key can read it.
func encryptTo(pub *rsa.PublicKey, secret []byte, label string) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, secret, []byte(label))
}

// checkRequest returns nil only if sig is pub's RSA-PSS signature over req's SHA-256 digest.
func checkRequest(pub *rsa.PublicKey, req, sig []byte) error {
	digest := sha256.Sum256(req)
	return rsa.VerifyPSS(pub, crypto.SHA256, digest[:], sig, nil)
}
