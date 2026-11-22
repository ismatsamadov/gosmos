package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// printRequest prints the request line and the headers that req writes, one per line.
func printRequest(req *http.Request) {
	var wire bytes.Buffer
	if err := req.Write(&wire); err != nil {
		fmt.Println("error:", err)
		return
	}
	head, _, _ := strings.Cut(wire.String(), "\r\n\r\n")
	fmt.Println(strings.ReplaceAll(head, "\r\n", "\n"))
}
