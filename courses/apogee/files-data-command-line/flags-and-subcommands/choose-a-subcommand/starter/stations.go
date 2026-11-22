package main

import (
	"fmt"
	"io"
)

// Station is one of Apogee's ground stations.
type Station struct {
	Code string // "YST": three capital letters
	Name string // "Ystmark"
}

// stations is Apogee's station table, in the order the team lists it.
var stations = []Station{
	{Code: "YST", Name: "Ystmark"},
	{Code: "PEL", Name: "Port Ellery"},
	{Code: "SQT", Name: "Salar Quieto"},
	{Code: "LAN", Name: "Lantern Atoll"},
	{Code: "CWR", Name: "Cape Wren"},
}

// runStations prints the station table, one code and name per line.
func runStations(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("stations", "stations")
	if status, ok := parseFlags(fs, args, stdout, stderr); !ok {
		return status
	}
	if fs.NArg() > 0 {
		return usageError(fs, stderr, "stations takes no arguments")
	}
	for _, s := range stations {
		fmt.Fprintln(stdout, s.Code, s.Name)
	}
	return 0
}
