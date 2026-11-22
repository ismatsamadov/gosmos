package main

// canonicalCase is one command and the text canonical must make of it.
type canonicalCase struct {
	name string
	cmd  Command
	want string
}

// canonicalCases are the cases the tests run against canonical.
var canonicalCases = []canonicalCase{
	{
		name: "one argument",
		cmd:  Command{SatID: "pipit-2", Seq: 118, Op: "SET_MODE", Args: []string{"SAFE"}},
		want: "pipit-2|118|SET_MODE|SAFE",
	},
	{
		name: "no arguments",
		cmd:  Command{SatID: "pipit-1", Seq: 7, Op: "PING"},
		want: "pipit-1|7|PING",
	},
	{
		name: "two arguments",
		cmd:  Command{SatID: "furrow-a", Seq: 5531, Op: "DUMP_FILE", Args: []string{"img-0511", "0"}},
		want: "furrow-a|5531|DUMP_FILE|img-0511|0",
	},
	{
		name: "signature left out",
		cmd:  Command{SatID: "cirrine-1", Seq: 90, Op: "REBOOT", Sig: "9f2c"},
		want: "cirrine-1|90|REBOOT",
	},
}
