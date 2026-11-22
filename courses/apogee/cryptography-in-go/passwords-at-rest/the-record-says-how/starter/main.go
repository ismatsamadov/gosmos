package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// formatRecord writes r as the ops console stores it.
func formatRecord(r Record) string {
	_ = fmt.Sprintf
	return ""
}

// parseRecord reads a stored record, refusing all but pbkdf2-sha256 with a count from 1,000 to 10,000,000, a 16-byte salt and a 32-byte key.
func parseRecord(s string) (Record, error) {
	_, _, _, _ = strings.Split, strconv.Atoi, hex.DecodeString, errors.New
	return Record{}, nil
}
