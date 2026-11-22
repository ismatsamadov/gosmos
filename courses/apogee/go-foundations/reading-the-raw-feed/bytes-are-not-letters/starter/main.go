package main

import "fmt"

func main() {
	names := []string{"Ystmark", "Quebrada Añil", "Tomás"}
	for _, name := range names {
		fmt.Printf("%s: %d\n", name, len(name))
	}
}
