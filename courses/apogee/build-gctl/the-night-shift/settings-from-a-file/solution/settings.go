package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"
)

// Settings are what gctl.json can change. Every field is optional.
type Settings struct {
	Prompt   string  `json:"prompt"`    // "gctl> " unless set
	CacheTTL ttlText `json:"cache_ttl"` // how long passes stay fresh: 10 minutes unless set
	Station  string  `json:"station"`   // the station passes shows with no argument; "" for all
}

// defaultSettings returns the settings gctl uses wherever gctl.json says nothing.
func defaultSettings() Settings {
	return Settings{Prompt: "gctl> ", CacheTTL: ttlText(10 * time.Minute)}
}

// loadSettings reads the settings at path strictly; a missing file gives the defaults.
func loadSettings(path string) (Settings, error) {
	st := defaultSettings()
	f, err := os.Open(path)
	if errors.Is(err, fs.ErrNotExist) {
		return st, nil
	}
	if err != nil {
		return Settings{}, err
	}
	defer f.Close()
	if err := decodeStrict(f, &st); err != nil {
		return Settings{}, fmt.Errorf("%s: %w", path, err)
	}
	return st, nil
}

// apply puts settings into the shell: its prompt, its cache's time to live and its default
// station.
func (s *Shell) apply(st Settings) {
	s.prompt = st.Prompt
	s.cache.ttl = time.Duration(st.CacheTTL)
	s.station = st.Station
}

// decodeStrict decodes the one JSON value in r into v, refusing unknown fields and anything
// after the value.
func decodeStrict(r io.Reader, v any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return errors.New("more data after the settings")
	}
	return nil
}

// ttlText is a time to live that gctl.json writes as duration text, such as "10m" or "90s".
type ttlText time.Duration

// UnmarshalJSON reads duration text, and refuses anything but a positive duration.
func (d *ttlText) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return fmt.Errorf("cache_ttl %s is not duration text", data)
	}
	v, err := time.ParseDuration(text)
	if err != nil || v <= 0 {
		return fmt.Errorf("cache_ttl %q is not a positive duration", text)
	}
	*d = ttlText(v)
	return nil
}

// String writes the time to live as duration text.
func (d ttlText) String() string {
	return time.Duration(d).String()
}
