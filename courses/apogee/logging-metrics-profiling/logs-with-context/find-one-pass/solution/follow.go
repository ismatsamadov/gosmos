package main

import (
	"bufio"
	"encoding/json"
	"io"
)

// follow reads JSON log lines from r and returns pass passID's lines, in the order read, and how many lines weren't JSON.
func follow(r io.Reader, passID string) (lines []Line, broken int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		var l Line
		if err := json.Unmarshal(sc.Bytes(), &l); err != nil {
			broken++
			continue
		}
		if l.PassID == passID {
			l.Raw = sc.Text()
			lines = append(lines, l)
		}
	}
	return lines, broken, sc.Err()
}
