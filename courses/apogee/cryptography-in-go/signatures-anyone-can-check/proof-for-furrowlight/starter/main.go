package main

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

// checkDelivery checks a delivered file against the envelope of its receipt, with the receipt
// keys Apogee publishes, and returns the receipt.
func checkDelivery(keys []ReceiptKey, env []byte, file io.Reader) (Receipt, error) {
	_ = time.RFC3339
	return Receipt{}, fmt.Errorf("delivery not checked")
}

// checkFile hashes file with the receipt's own hash and returns ErrWrongFile unless its size
// and digest are the receipt's.
func checkFile(r Receipt, file io.Reader) error {
	_ = bytes.Equal
	return fmt.Errorf("file not checked")
}
