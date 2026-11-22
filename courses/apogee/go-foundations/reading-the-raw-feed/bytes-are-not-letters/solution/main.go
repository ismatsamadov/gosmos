package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	names := []string{"Ystmark", "Quebrada Añil", "Tomás"}
	for _, name := range names {
		fmt.Printf("%s: %d\n", name, utf8.RuneCountInString(name))
	}
}
