package main

import "crypto/cipher"

// sealWithHeader seals plaintext under aead with a fresh nonce and returns s with its nonce and
// ciphertext; s's header stays readable.
func sealWithHeader(aead cipher.AEAD, s Sealed, plaintext []byte) Sealed {
	s.Nonce = newNonce()
	s.Ciphertext = aead.Seal(nil, s.Nonce, plaintext, nil)
	return s
}

// openWithHeader opens s's ciphertext under aead.
func openWithHeader(aead cipher.AEAD, s Sealed) ([]byte, error) {
	return aead.Open(nil, s.Nonce, s.Ciphertext, nil)
}
