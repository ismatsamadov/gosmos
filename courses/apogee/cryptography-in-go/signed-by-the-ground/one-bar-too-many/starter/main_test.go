package main

import (
	"fmt"
	"strings"
	"testing"
)

// cmd returns Pipit-2's command seq with op and args.
func cmd(seq uint64, op string, args ...string) Command {
	return Command{SatID: "pipit-2", Seq: seq, Op: op, Args: args}
}

// show writes c for a message, its arguments quoted so spaces and empty ones show.
func show(c Command) string {
	return fmt.Sprintf("{%q %d %q %q}", c.SatID, c.Seq, c.Op, c.Args)
}

func TestCanonical(t *testing.T) {
	tests := []struct {
		name string
		c    Command
		want string
	}{
		{"the lab bench's command", cmd(118, "SET_MODE", "SAFE"), "pipit-2|118|SET_MODE|SAFE"},
		{"two arguments", cmd(118, "SET_MODE", "SAFE", "REBOOT"), "pipit-2|118|SET_MODE|SAFE|REBOOT"},
		{"no arguments", cmd(119, "PING"), "pipit-2|119|PING"},
		{"a file to dump", Command{SatID: "cirrine-1", Seq: 4471, Op: "DUMP_FILE", Args: []string{"occ-0825.bin"}}, "cirrine-1|4471|DUMP_FILE|occ-0825.bin"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := canonical(tt.c)
			if err != nil || got != tt.want {
				t.Errorf("canonical(%s) = %q, %v, want %q, nil", show(tt.c), got, err, tt.want)
			}
		})
	}
}

func TestCanonicalRefuses(t *testing.T) {
	tests := []struct {
		name  string
		c     Command
		field string // what the error must start with
	}{
		{"a bar inside an argument", cmd(118, "SET_MODE", "SAFE|REBOOT"), `argument 1 contains "|"`},
		{"an empty argument", cmd(118, "SET_MODE", "SAFE", ""), "argument 2"},
		{"a space inside an argument", cmd(118, "SET_MODE", "SAFE REBOOT"), "argument 1"},
		{"a tab inside an argument", cmd(118, "SET_MODE", "SAFE\tREBOOT"), "argument 1"},
		{"sequence number 0", cmd(0, "SET_MODE", "SAFE"), "sequence number"},
		{"the satellite in capitals", Command{SatID: "Pipit-2", Seq: 118, Op: "SET_MODE", Args: []string{"SAFE"}}, "satellite ID"},
		{"no satellite", Command{Seq: 118, Op: "PING"}, "satellite ID"},
		{"the operation in lowercase", cmd(118, "set_mode", "SAFE"), "operation"},
		{"no operation", cmd(118, ""), "operation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := canonical(tt.c)
			if err == nil || !strings.HasPrefix(err.Error(), tt.field) {
				t.Errorf("canonical(%s) = %q, %v, want an error starting %#q", show(tt.c), got, err, tt.field)
			}
		})
	}
}
