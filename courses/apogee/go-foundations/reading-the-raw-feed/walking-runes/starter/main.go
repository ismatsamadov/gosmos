package main

import (
	"fmt"
	"strings"
)

// callSign turns a name into a radio call sign: only A to Z, digits, hyphens and underscores.
func callSign(name string) string {
	var b strings.Builder
	// TODO: write each rune of name to b in its call-sign form.
	return b.String()
}

func main() {
	fmt.Println(callSign("Salar Quieto"))
}
