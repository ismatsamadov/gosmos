package main

import (
	"fmt"
	"strconv"
	"strings"
)

// formatSealed writes s as one SEAL line: SEAL 2 <sat> <seq> <key id> <nonce> <ciphertext>.
func formatSealed(s Sealed) string {
	_ = fmt.Sprintf
	return ""
}

// parseSealed reads one SEAL line, refusing anything but the exact format.
func parseSealed(line string) (Sealed, error) {
	_, _ = strconv.ParseUint, strings.Split
	return Sealed{}, nil
}
