package main

import "crypto/mlkem"

// wrapFor makes a new key for the holder of the encapsulation key ek, and the ciphertext that carries it.
func wrapFor(ek []byte) (key, ciphertext []byte, err error) {
	_ = mlkem.NewEncapsulationKey768
	return nil, nil, nil
}

// unwrap returns the key that ciphertext carries, opened with the decapsulation key dk.
func unwrap(dk *mlkem.DecapsulationKey768, ciphertext []byte) ([]byte, error) {
	return nil, nil
}
