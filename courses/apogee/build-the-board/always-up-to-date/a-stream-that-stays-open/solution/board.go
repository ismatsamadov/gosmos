package main

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"sync"
	"time"
)

// Sources is where the Board's data comes from.
type Sources struct {
	Passes  func(ctx context.Context) ([]Booking, error)            // passplan's schedule: every booked pass, in any order
	Weather func(ctx context.Context, code string) (Weather, error) // the weather service: WeatherClient.Reading
	Rota    []Shift                                                 // the on-call rota, oldest first
}

// Shift is one week on call; it starts on a Monday at 08:00 UTC.
type Shift struct {
	From time.Time
	Name string // "Nkem Halloran"
}

// View is everything the Board shows at one moment.
type View struct {
	Passes  []Booking `json:"passes"`  // LOS after now, by AOS then ID, at most 12; [] when none
	Weather []Panel   `json:"weather"` // one per station, in the station table's order
	OnCall  string    `json:"on_call"` // the latest shift that started at or before now
	version int       // which view this is: Refresh numbers them from 1
}

// Panel is one station's weather on the Board.
type Panel struct {
	Station string   `json:"station"`           // "CWR"
	Weather *Weather `json:"weather,omitempty"` // nil when there is no reading
	Call    string   `json:"call"`              // "track", "warn", "stow", or "none" without a reading
}

// Board is the program behind the Board: its sources, its page and what it shows now.
type Board struct {
	src      Sources
	stations []Station
	now      func() time.Time
	page     *template.Template // board.html, parsed once by NewBoard

	mu      sync.Mutex // guards view, json and viewers
	view    View       // what the Board shows now: an empty view, version 0, until Refresh
	json    []byte     // view as JSON
	viewers int        // streams open now, for /health
}

// NewBoard returns a Board that reads src, shows stations in their order and tells the
// time with now. It shows nothing until its first Refresh.
func NewBoard(src Sources, stations []Station, now func() time.Time) (*Board, error) {
	page, err := parsePage()
	if err != nil {
		return nil, fmt.Errorf("board.html: %w", err)
	}
	return &Board{src: src, stations: slices.Clone(stations), now: now, page: page}, nil
}

// Refresh builds a new view from the sources at b.now() and shows it. When passplan
// fails, the Board keeps the view it had and Refresh returns the error.
func (b *Board) Refresh(ctx context.Context) error {
	now := b.now()
	all, err := b.src.Passes(ctx)
	if err != nil {
		return fmt.Errorf("passes: %w", err)
	}
	v := View{Passes: nextPasses(all, now), Weather: b.panels(ctx), OnCall: onCall(b.src.Rota, now)}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if bytes.Equal(data, b.json) {
		return nil
	}
	v.version = b.view.version + 1
	b.view, b.json = v, data
	return nil
}

// Handler returns the Board's routes.
func (b *Board) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", b.servePage)
	mux.HandleFunc("GET /board.json", b.serveJSON)
	mux.HandleFunc("GET /events", b.serveEvents)
	mux.HandleFunc("GET /health", b.serveHealth)
	return mux
}

// current returns the view the Board shows now.
func (b *Board) current() View {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.view
}

// panels asks the weather service for each station's reading, in the stations' order.
// A station whose reading fails gets a panel without one.
func (b *Board) panels(ctx context.Context) []Panel {
	panels := make([]Panel, 0, len(b.stations))
	for _, st := range b.stations {
		p := Panel{Station: st.Code, Call: "none"}
		if w, err := b.src.Weather(ctx, st.Code); err == nil {
			p.Weather, p.Call = &w, callFor(w.WindKmh)
		}
		panels = append(panels, p)
	}
	return panels
}

// callFor is Cape Wren's rule for a dish in a sustained wind of windKmh: track under
// 50 km/h, warn from 50, stow from 72.
func callFor(windKmh float64) string {
	switch {
	case windKmh >= 72:
		return "stow"
	case windKmh >= 50:
		return "warn"
	default:
		return "track"
	}
}

// onCall returns the name on the latest shift of rota that started at or before now, or
// "" before the first.
func onCall(rota []Shift, now time.Time) string {
	name := ""
	for _, s := range rota {
		if s.From.After(now) {
			break
		}
		name = s.Name
	}
	return name
}

// nextPasses returns, in a new slice, the passes of all whose LOS is after now, by AOS
// and then booking ID, at most 12: the rows the Board's television fits.
func nextPasses(all []Booking, now time.Time) []Booking {
	next := make([]Booking, 0, len(all))
	for _, p := range all {
		if p.LOS.After(now) {
			next = append(next, p)
		}
	}
	slices.SortFunc(next, func(a, b Booking) int {
		return cmp.Or(a.AOS.Compare(b.AOS), cmp.Compare(a.ID, b.ID))
	})
	return next[:min(len(next), 12)]
}
