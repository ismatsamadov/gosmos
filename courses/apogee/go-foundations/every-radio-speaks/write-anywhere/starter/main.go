package main

import (
	"fmt"
	"io"
	"os"
)

// writeBoard writes the Board's header to w, then one pass ID per line.
func writeBoard(w io.Writer, ids []string) {
	fmt.Println("NEXT PASSES")
	for _, id := range ids {
		fmt.Println(id)
	}
}

func main() {
	writeBoard(os.Stdout, []string{"pipit-1/PEL/04182", "furrow-a/SQT/11873", "thistle/LAN/02218"})
}
