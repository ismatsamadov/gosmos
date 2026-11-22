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

// loadBookings reads the bookings in the named file, or on stdin when no file is named.
func loadBookings(name string, stdin io.Reader) ([]Booking, error) {
	if name == "" {
		return readBookings(stdin)
	}
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

// runPasses lists the bookings in the named file, or on stdin, one per line by AOS.
func runPasses(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := newFlagSet("passes", "passes [-station CODE] [file]")
	station := fs.String("station", "", "show only the bookings at station `CODE`")
	if status, ok := parseFlags(fs, args, stdout, stderr); !ok {
		return status
	}
	if fs.NArg() > 1 {
		return usageError(fs, stderr, "passes reads one file at most")
	}
	if _, ok := findStation(*station); *station != "" && !ok {
		return report(nil, fmt.Errorf("unknown station %q", *station), stdout, stderr)
	}
	bookings, err := loadBookings(fs.Arg(0), stdin)
	return report(passLines(bookings, *station), err, stdout, stderr)
}
