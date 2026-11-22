package main

import "crypto/mlkem"

// wrapFor makes a new key for the holder of the encapsulation key ek, and the ciphertext that carries it.
func wrapFor(ek []byte) (key, ciphertext []byte, err error) {
	pub, err := mlkem.NewEncapsulationKey768(ek)
	if err != nil {
		return nil, nil, err
	}
	key, ciphertext = pub.Encapsulate()
	return key, ciphertext, nil
}

// unwrap returns the key that ciphertext carries, opened with the decapsulation key dk.
func unwrap(dk *mlkem.DecapsulationKey768, ciphertext []byte) ([]byte, error) {
	return dk.Decapsulate(ciphertext)
}
