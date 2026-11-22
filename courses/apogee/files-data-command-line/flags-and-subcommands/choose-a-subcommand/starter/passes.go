package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
)

// Booking is one reserved pass, as the booking export writes it.
type Booking struct {
	ID      string    `json:"id"`
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"`
}

// readBookings decodes a booking log, one JSON booking per line, returning those read before any error.
func readBookings(r io.Reader) ([]Booking, error) {
	dec := json.NewDecoder(r)
	var bookings []Booking
	for {
		var b Booking
		err := dec.Decode(&b)
		if err == io.EOF {
			return bookings, nil
		}
		if err != nil {
			return bookings, fmt.Errorf("booking %d: %w", len(bookings)+1, err)
		}
		bookings = append(bookings, b)
	}
}

// readFile reads the bookings in the named file.
func readFile(name string) ([]Booking, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readBookings(f)
}

// passLines formats the bookings at station, or at every station when it is empty, one line each in order of AOS.
func passLines(bookings []Booking, station string) []string {
	var picked []Booking
	for _, b := range bookings {
		if station == "" || b.Station == station {
			picked = append(picked, b)
		}
	}
	slices.SortStableFunc(picked, func(a, b Booking) int { return a.AOS.Compare(b.AOS) })
	var lines []string
	for _, b := range picked {
		lines = append(lines, fmt.Sprintf("%s %s %s %s", b.AOS.UTC().Format(time.RFC3339), b.Station, b.SatID, b.ID))
	}
	return lines
}

// runPasses lists the bookings in a file, one per line by AOS.
func runPasses(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("passes", "passes [-station CODE] file")
	station := fs.String("station", "", "show only the bookings at station `CODE`")
	if status, ok := parseFlags(fs, args, stdout, stderr); !ok {
		return status
	}
	if fs.NArg() != 1 {
		return usageError(fs, stderr, "passes needs one file")
	}
	bookings, err := readFile(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "gctl: %v\n", err)
		return 1
	}
	for _, line := range passLines(bookings, *station) {
		fmt.Fprintln(stdout, line)
	}
	return 0
}
