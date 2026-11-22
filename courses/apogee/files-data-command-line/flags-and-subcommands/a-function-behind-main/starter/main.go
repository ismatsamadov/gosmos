package main

import (
	"fmt"
	"io"
	"os"
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

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run prints the station table, one code and name per line, and returns gctl's exit status.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	_ = fmt.Fprintln
	return 0
}
