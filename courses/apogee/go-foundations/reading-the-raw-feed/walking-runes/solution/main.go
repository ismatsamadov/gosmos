package main

import (
	"fmt"
	"strings"
	"unicode"
)

// callSign turns a name into a radio call sign: only A to Z, digits, hyphens and underscores.
func callSign(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case 'a' <= r && r <= 'z':
			b.WriteRune(unicode.ToUpper(r))
		case 'A' <= r && r <= 'Z', '0' <= r && r <= '9', r == '-':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune('-')
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

func main() {
	fmt.Println(callSign("Salar Quieto"))
}
