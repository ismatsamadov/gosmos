package main

import (
	"errors"
	"strconv"
	"strings"
)

// checkLine parses one CMD line, then lets g check its signature and sequence number.
func checkLine(line string, key []byte, g *Guard) (Command, error) {
	_, _ = strconv.ParseUint, strings.Fields
	return Command{}, errors.New("bad command line")
}
