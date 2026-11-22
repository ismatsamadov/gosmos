package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// shortCode returns the first n bytes of text's SHA-256 in hex, as Brinewatch's tool shows a
// fingerprint.
func shortCode(text string, n int) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:n])
}

// firstMatch tries the texts "brinewatch 0", "brinewatch 1", ... in order, at most limit of
// them, and returns the numbers of the first two whose codes of n bytes are equal.
func firstMatch(n, limit int) (a, b int, ok bool) {
	seen := map[string]int{}
	for i := range limit {
		code := shortCode(fmt.Sprintf("brinewatch %d", i), n)
		if j, found := seen[code]; found {
			return j, i, true
		}
		seen[code] = i
	}
	return 0, 0, false
}

func main() {
	for _, n := range []int{2, 3} {
		a, b, ok := firstMatch(n, 100000)
		if !ok {
			fmt.Printf("%d bytes: no two match in 100000 tries\n", n)
			continue
		}
		code := shortCode(fmt.Sprintf("brinewatch %d", b), n)
		fmt.Printf("%d bytes: %d and %d share %s, %d tries\n", n, a, b, code, b+1)
	}
}
