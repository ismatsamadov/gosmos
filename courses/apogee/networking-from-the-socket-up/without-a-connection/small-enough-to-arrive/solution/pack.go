package main

import "fmt"

// pack puts lines, in order, into datagrams for the live port: whole lines
// joined by "\n", none over maxDatagram bytes.
func pack(lines []string) ([][]byte, error) {
	var out [][]byte
	var d []byte
	for _, line := range lines {
		if len(line) > maxDatagram {
			return nil, fmt.Errorf("line of %d bytes: %w", len(line), ErrTooBig)
		}
		if len(d) > 0 && len(d)+1+len(line) > maxDatagram {
			out = append(out, d)
			d = nil
		}
		if len(d) > 0 {
			d = append(d, '\n')
		}
		d = append(d, line...)
	}
	if len(d) > 0 {
		out = append(out, d)
	}
	return out, nil
}
