package main

import (
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
)

// handOver checks Brinewatch's request, then returns the two lines that carry its new secrets,
// each locked to the receiving key the request names.
func handOver(signKey *rsa.PublicKey, r Request, apiKey, webhook []byte) ([]string, error) {
	if err := checkRequest(signKey, r.Text, r.Sig); err != nil {
		return nil, fmt.Errorf("request not signed by Brinewatch: %w", err)
	}
	pub, err := readRSA(r.Key)
	if err != nil {
		return nil, err
	}
	fp, err := keyFingerprint(pub)
	if err != nil {
		return nil, err
	}
	if string(r.Text) != "brinewatch handover "+fp {
		return nil, errors.New("the request names another key")
	}
	lockedKey, err := encryptTo(pub, apiKey, apiKeyLabel)
	if err != nil {
		return nil, err
	}
	lockedHook, err := encryptTo(pub, webhook, webhookLabel)
	if err != nil {
		return nil, err
	}
	return []string{
		"api-key " + hex.EncodeToString(lockedKey),
		"webhook-secret " + hex.EncodeToString(lockedHook),
	}, nil
}
