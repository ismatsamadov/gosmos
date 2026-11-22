package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Command is one telecommand on its way up.
type Command struct {
	SatID string
	Seq   uint64 // strictly increasing per satellite
	Op    string // "PING", "SET_MODE", "DUMP_FILE", "REBOOT"
	Args  []string
	Sig   string // hex HMAC-SHA256 of the canonical form
}

// canonical lays a command out for signing: its fields joined by "|", arguments last, no Sig.
func canonical(c Command) string {
	fields := []string{c.SatID, strconv.FormatUint(c.Seq, 10), c.Op}
	fields = append(fields, c.Args...)
	return strings.Join(fields, "|")
}

func main() {
	cmd := Command{SatID: "pipit-2", Seq: 118, Op: "SET_MODE", Args: []string{"SAFE"}}
	fmt.Println(canonical(cmd))
}
