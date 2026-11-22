package main

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// Plan is a calibration day: the station and its booked passes that day, in AOS order.
type Plan struct {
	Day     string
	Station Station
	Passes  []Booking
}

// planDay reads the settings and the bookings in dir and plans the day.
func planDay(dir, format string) (Plan, error) {
	settings, err := readSettings(filepath.Join(dir, "settings.json"))
	if err != nil {
		return Plan{}, fmt.Errorf("settings: %w", err)
	}
	bookings, err := readBookingFile(filepath.Join(dir, "bookings."+format), format)
	if err != nil {
		return Plan{}, fmt.Errorf("bookings: %w", err)
	}
	stations, err := loadStations()
	if err != nil {
		return Plan{}, err
	}
	if err := validateBookings(bookings, stations); err != nil {
		return Plan{}, fmt.Errorf("bookings: %w", err)
	}
	return newPlan(settings, stations, bookings)
}

// readSettings opens a settings file and loads it.
func readSettings(path string) (Settings, error) {
	file, err := os.Open(path)
	if err != nil {
		return Settings{}, err
	}
	defer file.Close()
	return loadConfig(file)
}

// readBookingFile opens a booking file and reads it in the given format.
func readBookingFile(path, format string) ([]Booking, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return readBookings(file, format)
}

// validateBookings checks every booking, returning the first problem it finds.
func validateBookings(bookings []Booking, stations map[string]Station) error {
	for i, b := range bookings {
		_, known := stations[b.Station]
		switch {
		case b.ID == "":
			return fmt.Errorf("booking %d: no id", i+1)
		case b.AOS.IsZero():
			return fmt.Errorf("booking %s: no AOS", b.ID)
		case !b.LOS.After(b.AOS):
			return fmt.Errorf("booking %s: LOS is not after AOS", b.ID)
		case !known:
			return fmt.Errorf("booking %s: unknown station %q", b.ID, b.Station)
		}
	}
	return nil
}

// newPlan keeps the bookings at the settings' station on the settings' day, in order.
func newPlan(settings Settings, stations map[string]Station, bookings []Booking) (Plan, error) {
	station, ok := stations[settings.Station]
	if !ok {
		return Plan{}, fmt.Errorf("settings: unknown station %q", settings.Station)
	}
	var passes []Booking
	for _, b := range bookings {
		if b.Station == settings.Station && b.AOS.UTC().Format(time.DateOnly) == settings.Day {
			passes = append(passes, b)
		}
	}
	return Plan{Day: settings.Day, Station: station, Passes: orderedPasses(passes)}, nil
}

// orderedPasses returns a copy of passes sorted by AOS, and by booking ID when two start together.
func orderedPasses(passes []Booking) []Booking {
	sorted := slices.Clone(passes)
	slices.SortFunc(sorted, func(a, b Booking) int {
		return cmp.Or(a.AOS.Compare(b.AOS), strings.Compare(a.ID, b.ID))
	})
	return sorted
}

// writeReport writes the day report: a heading, the number of booked passes, then one line each.
func writeReport(w io.Writer, p Plan) error {
	var b strings.Builder
	fmt.Fprintf(&b, "Calibration day %s at %s (%s)\n", p.Day, p.Station.Code, p.Station.Name)
	fmt.Fprintf(&b, "Booked passes (UTC): %d\n", len(p.Passes))
	for _, pass := range p.Passes {
		fmt.Fprintf(&b, "%s-%s %s %s", pass.AOS.UTC().Format(time.TimeOnly), pass.LOS.UTC().Format(time.TimeOnly), pass.SatID, pass.ID)
		if note := strings.Join(strings.Fields(pass.Note), " "); note != "" {
			b.WriteString(" " + note)
		}
		b.WriteString("\n")
	}
	_, err := io.WriteString(w, b.String())
	return err
}
