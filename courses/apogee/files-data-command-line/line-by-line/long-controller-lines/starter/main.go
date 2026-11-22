package main

import (
	"bufio"
	"io"
)

// maxRecordBytes is the most the scanner may hold for one record: 256 KiB.
const maxRecordBytes = 256 * 1024

// scanRecords counts the records in r, letting each be up to maxRecordBytes long.
func scanRecords(r io.Reader) (int, error) {
	sc := bufio.NewScanner(r)
	// TODO: let the scanner's buffer grow to maxRecordBytes.
	n := 0
	for sc.Scan() {
		n++
	}
	return n, sc.Err()
}
