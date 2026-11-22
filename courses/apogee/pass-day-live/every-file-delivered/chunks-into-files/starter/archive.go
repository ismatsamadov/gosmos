package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// loadArchive reads passes.txt in dir: passplan's pass archive for the window.
func loadArchive(dir string) ([]PassView, error) {
	f, err := os.Open(filepath.Join(dir, "passes.txt"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readArchive(f)
}

// readArchive reads a pass archive, one pass per line. A line it can't read is an error
// that names it: "archive line 3: want 4 fields, got 3".
func readArchive(r io.Reader) ([]PassView, error) {
	var passes []PassView
	sc := bufio.NewScanner(r)
	for n := 1; sc.Scan(); n++ {
		p, err := parsePass(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("archive line %d: %w", n, err)
		}
		passes = append(passes, p)
	}
	return passes, sc.Err()
}

// parsePass reads one line of the archive: the pass ID, AOS, LOS and maximum elevation,
// as "thistle-2/PEL/02136 2029-07-26T16:03:46Z 2029-07-26T16:09:59Z 11.7".
func parsePass(line string) (PassView, error) {
	f := strings.Fields(line)
	if len(f) != 4 {
		return PassView{}, fmt.Errorf("want 4 fields, got %d", len(f))
	}
	id := strings.Split(f[0], "/")
	if len(id) != 3 {
		return PassView{}, fmt.Errorf("pass ID %q is not satellite/station/orbit", f[0])
	}
	aos, err := time.Parse(time.RFC3339, f[1])
	if err != nil {
		return PassView{}, fmt.Errorf("AOS %q is not a time like 2029-07-26T16:03:46Z", f[1])
	}
	los, err := time.Parse(time.RFC3339, f[2])
	if err != nil {
		return PassView{}, fmt.Errorf("LOS %q is not a time like 2029-07-26T16:09:59Z", f[2])
	}
	elev, err := strconv.ParseFloat(f[3], 64)
	if err != nil {
		return PassView{}, fmt.Errorf("maximum elevation %q is not a number", f[3])
	}
	return PassView{ID: f[0], Pass: Pass{SatID: id[0], Station: id[1], AOS: aos, LOS: los, MaxElevDeg: elev}}, nil
}
