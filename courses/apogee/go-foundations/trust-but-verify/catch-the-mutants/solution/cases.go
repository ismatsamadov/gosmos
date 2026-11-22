package main

// acceptCase is one call to accept and the answer it must give.
type acceptCase struct {
	name      string
	seq, last uint64
	want      bool
}

// acceptCases are the cases the tests run against accept and against every mutant.
var acceptCases = []acceptCase{
	{"far ahead", 500, 118, true},
	{"next command", 119, 118, true},
	{"replay", 118, 118, false},
	{"old command", 117, 118, false},
}
