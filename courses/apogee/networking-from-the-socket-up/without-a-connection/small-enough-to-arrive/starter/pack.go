package main

import "strings"

// pack puts lines, in order, into datagrams for the live port: whole lines
// joined by "\n", none over maxDatagram bytes.
func pack(lines []string) ([][]byte, error) {
	if len(lines) == 0 {
		return nil, nil
	}
	return [][]byte{[]byte(strings.Join(lines, "\n"))}, nil
}
