package main

import (
	"errors"
	"strconv"
	"strings"
)

// checkLine parses one CMD line, then lets g check its signature and sequence number.
func checkLine(line string, key []byte, g *Guard) (Command, error) {
	c, ok := parseCmdLine(line)
	if !ok {
		return Command{}, errors.New("bad command line")
	}
	if err := g.Accept(c, key); err != nil {
		return Command{}, err
	}
	return c, nil
}

// parseCmdLine reads CMD, a satellite, a sequence number, an operation, any arguments and one
// last sig= field of 64 lowercase hex digits, and reports whether the line was exactly that.
func parseCmdLine(line string) (Command, bool) {
	f := strings.Fields(line)
	if len(f) < 5 || f[0] != "CMD" {
		return Command{}, false
	}
	sig, ok := strings.CutPrefix(f[len(f)-1], "sig=")
	seq, err := strconv.ParseUint(f[2], 10, 64)
	if !ok || len(sig) != 64 || strings.Trim(sig, "0123456789abcdef") != "" || err != nil {
		return Command{}, false
	}
	args := f[4 : len(f)-1]
	for _, arg := range args {
		if strings.HasPrefix(arg, "sig=") {
			return Command{}, false
		}
	}
	return Command{SatID: f[1], Seq: seq, Op: f[3], Args: args, Sig: sig}, true
}
