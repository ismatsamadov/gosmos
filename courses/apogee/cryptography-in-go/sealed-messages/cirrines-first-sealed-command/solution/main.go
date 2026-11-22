package main

import "strings"

// sealCommand seals c for its satellite under master version keyID and returns its SEAL line.
func sealCommand(ring *KeyRing, keyID string, c Command) (string, error) {
	if _, err := canonical(c); err != nil {
		return "", err
	}
	key, err := ring.SatKey(keyID, "seal", c.SatID)
	if err != nil {
		return "", err
	}
	aead, err := newAEAD(key)
	if err != nil {
		return "", err
	}
	s := Sealed{Version: 2, SatID: c.SatID, Seq: c.Seq, KeyID: keyID, Nonce: newNonce()}
	plaintext := strings.Join(append([]string{c.Op}, c.Args...), "|")
	s.Ciphertext = aead.Seal(nil, s.Nonce, []byte(plaintext), s.AD())
	return formatSealed(s), nil
}
