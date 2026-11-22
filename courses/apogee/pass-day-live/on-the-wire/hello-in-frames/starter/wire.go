package main

import (
	"bufio"
	"fmt"
	"strings"
)

// readBanner reads the station's banner, "<CODE> READY PROTO 1" ended by CR LF or LF, and
// refuses any other line.
func readBanner(br *bufio.Reader, code string) error {
	line, err := br.ReadString('\n')
	if err != nil {
		return err
	}
	text := strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
	if text != code+" READY PROTO 1" {
		return fmt.Errorf("bad banner %q", text)
	}
	return nil
}
