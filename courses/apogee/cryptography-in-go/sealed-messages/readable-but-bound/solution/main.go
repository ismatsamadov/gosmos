package main

import "crypto/cipher"

// sealWithHeader seals plaintext under aead with a fresh nonce and returns s with its nonce and
// ciphertext; s's header stays readable, and the tag covers it.
func sealWithHeader(aead cipher.AEAD, s Sealed, plaintext []byte) Sealed {
	s.Nonce = newNonce()
	s.Ciphertext = aead.Seal(nil, s.Nonce, plaintext, s.AD())
	return s
}

// openWithHeader opens s's ciphertext under aead, and only under the header it was sealed with.
func openWithHeader(aead cipher.AEAD, s Sealed) ([]byte, error) {
	return aead.Open(nil, s.Nonce, s.Ciphertext, s.AD())
}
