package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

// Booking is one reserved pass, as the booking export writes it.
type Booking struct {
	ID      string    `json:"id"`
	SatID   string    `json:"sat_id"`
	Station string    `json:"station"`
	AOS     time.Time `json:"aos"`
	LOS     time.Time `json:"los"`
	Note    string    `json:"note,omitempty"`
}

// readBookings reads bookings in the named format: json for the export, csv for the sheet.
func readBookings(r io.Reader, format string) ([]Booking, error) {
	switch format {
	case "json":
		return decodeBookings(r)
	case "csv":
		return readBookingSheet(r)
	default:
		return nil, fmt.Errorf("unsupported booking format %q", format)
	}
}

// decodeBookings reads a day's booking export: one JSON array of bookings.
func decodeBookings(r io.Reader) ([]Booking, error) {
	var bookings []Booking
	if err := json.NewDecoder(r).Decode(&bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// readBookingSheet reads a booking sheet: CSV whose header names its columns, in any order.
func readBookingSheet(r io.Reader) ([]Booking, error) {
	sheet := csv.NewReader(r)
	header, err := sheet.Read()
	if err != nil {
		return nil, fmt.Errorf("header: %w", err)
	}
	col := map[string]int{}
	for i, name := range header {
		if _, twice := col[name]; twice {
			return nil, fmt.Errorf("header: %s appears twice", name)
		}
		col[name] = i
	}
	for _, name := range []string{"ID", "SATELLITE", "STATION", "AOS", "LOS", "NOTE"} {
		if _, ok := col[name]; !ok {
			return nil, fmt.Errorf("header: no %s column", name)
		}
	}
	var bookings []Booking
	for {
		row, err := sheet.Read()
		if errors.Is(err, io.EOF) {
			return bookings, nil
		}
		if err != nil {
			return nil, err
		}
		b := Booking{ID: row[col["ID"]], SatID: row[col["SATELLITE"]], Station: row[col["STATION"]], Note: row[col["NOTE"]]}
		if b.AOS, err = time.Parse(time.RFC3339, row[col["AOS"]]); err != nil {
			return nil, fmt.Errorf("booking %s: AOS: %w", b.ID, err)
		}
		if b.LOS, err = time.Parse(time.RFC3339, row[col["LOS"]]); err != nil {
			return nil, fmt.Errorf("booking %s: LOS: %w", b.ID, err)
		}
		bookings = append(bookings, b)
	}
}
