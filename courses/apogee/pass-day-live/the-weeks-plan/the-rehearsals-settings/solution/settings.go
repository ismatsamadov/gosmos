package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"
)

// Settings is rehearsal.json: which week the rehearsal flies, and at which stations.
type Settings struct {
	Week     string   `json:"week"`     // the first day, "2029-07-23": the week starts at 00:00 UTC
	Days     int      `json:"days"`     // 1 to 7; 7 unless the file says
	Stations []string `json:"stations"` // station codes in any order; none means all twenty
}

// loadSettings reads rehearsal.json in dir strictly and checks it against the network.
func loadSettings(dir string, n *Network) (Settings, error) {
	f, err := os.Open(filepath.Join(dir, "rehearsal.json"))
	if err != nil {
		return Settings{}, err
	}
	defer f.Close()
	set := Settings{Days: 7}
	if err := decodeStrict(f, &set, "settings"); err != nil {
		return Settings{}, fmt.Errorf("rehearsal.json: %w", err)
	}
	if err := set.check(n); err != nil {
		return Settings{}, fmt.Errorf("rehearsal.json: %w", err)
	}
	return set, nil
}

// decodeStrict decodes the one JSON value in r into v, refusing unknown fields and
// anything after the value; what names the value in that error: "settings".
func decodeStrict(r io.Reader, v any, what string) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return errors.New("more data after the " + what)
	}
	return nil
}

// check refuses settings the rehearsal can't fly.
func (s Settings) check(n *Network) error {
	if s.Week == "" {
		return errors.New("week is required")
	}
	if _, err := time.Parse(time.DateOnly, s.Week); err != nil {
		return fmt.Errorf("week %q is not a date like 2006-01-02", s.Week)
	}
	if s.Days < 1 || s.Days > 7 {
		return fmt.Errorf("days %d is not from 1 to 7", s.Days)
	}
	for _, code := range s.Stations {
		if !slices.ContainsFunc(n.Sites, func(st Site) bool { return st.Code == code }) {
			return fmt.Errorf("unknown station %q", code)
		}
	}
	return nil
}

// sites returns the stations the rehearsal flies, in the network's table order: all of
// them when the settings name none.
func (s Settings) sites(n *Network) []Site {
	if len(s.Stations) == 0 {
		return n.Sites
	}
	var out []Site
	for _, st := range n.Sites {
		if slices.Contains(s.Stations, st.Code) {
			out = append(out, st)
		}
	}
	return out
}
