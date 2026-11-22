package main

import (
	"fmt"
	"strconv"
	"strings"
)

// canonical returns the text that gets signed for c: its fields joined by "|", arguments last.
func canonical(c Command) (string, error) {
	_ = fmt.Errorf
	fields := append([]string{c.SatID, strconv.FormatUint(c.Seq, 10), c.Op}, c.Args...)
	return strings.Join(fields, "|"), nil
}
