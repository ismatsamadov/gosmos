package main

// Holder is one certificate issueAll makes: its files' name, its common name, its serial, and
// whether it proves a client.
type Holder struct {
	Name   string // the files' name: "yst" for yst.pem and yst.key
	CN     string // "yst.apogee.internal", or a client's name: "uplinkd"
	Serial int64
	Client bool // uplinkd, passplan and dishd prove clients; the stations prove servers
}

// holders are the certificates of Friday 22 September 2028: five stations, then uplinkd, passplan
// and dishd.
var holders = []Holder{
	{Name: "yst", CN: "yst.apogee.internal", Serial: 1001},
	{Name: "pel", CN: "pel.apogee.internal", Serial: 1002},
	{Name: "sqt", CN: "sqt.apogee.internal", Serial: 1003},
	{Name: "lan", CN: "lan.apogee.internal", Serial: 1004},
	{Name: "cwr", CN: "cwr.apogee.internal", Serial: 1005},
	{Name: "uplinkd", CN: "uplinkd", Serial: 2001, Client: true},
	{Name: "passplan", CN: "passplan", Serial: 2002, Client: true},
	{Name: "dishd", CN: "dishd", Serial: 2003, Client: true},
}
