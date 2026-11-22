package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// night returns h:m:s on the night into Friday 11 February 2028, UTC.
func night(h, m, s int) time.Time {
	return time.Date(2028, time.February, 11, h, m, s, 0, time.UTC)
}

// tonight returns the night's bookings, as passplan sends them.
func tonight() []Booking {
	return []Booking{
		{ID: "bk-4986", SatID: "cirrine-1", Station: "PEL", AOS: night(3, 20, 0), LOS: night(3, 30, 0)},
		{ID: "bk-4988", SatID: "pipit-1", Station: "PEL", AOS: night(4, 12, 5), LOS: night(4, 21, 27), Note: "echo test"},
		{ID: "bk-4990", SatID: "pipit-2", Station: "YST", AOS: night(4, 31, 10), LOS: night(4, 40, 40)},
		{ID: "bk-4991", SatID: "furrow-a", Station: "CWR", AOS: night(4, 58, 30), LOS: night(5, 7, 13), Note: "<b>RUSH</b> & PRIORITY"},
		{ID: "bk-4992", SatID: "furrow-b", Station: "SQT", AOS: night(5, 12, 44), LOS: night(5, 22, 10), Note: "RESHOOT FIELD 12"},
		{ID: "bk-4993", SatID: "brine-3", Station: "LAN", AOS: night(5, 20, 15), LOS: night(5, 27, 0), Note: "harbor survey"},
		{ID: "bk-4994", SatID: "cirrine-1", Station: "YST", AOS: night(5, 47, 0), LOS: night(5, 57, 17)},
	}
}

// morning returns eight bookings for Friday morning, after the night's.
func morning() []Booking {
	return []Booking{
		{ID: "bk-4996", SatID: "pipit-2", Station: "YST", AOS: night(6, 5, 0), LOS: night(6, 14, 40)},
		{ID: "bk-4997", SatID: "brine-3", Station: "LAN", AOS: night(6, 41, 10), LOS: night(6, 49, 5)},
		{ID: "bk-4998", SatID: "pipit-1", Station: "PEL", AOS: night(7, 18, 40), LOS: night(7, 27, 5)},
		{ID: "bk-4999", SatID: "furrow-b", Station: "CWR", AOS: night(7, 55, 0), LOS: night(8, 3, 55)},
		{ID: "bk-5000", SatID: "cirrine-1", Station: "SQT", AOS: night(8, 30, 15), LOS: night(8, 41, 0)},
		{ID: "bk-5001", SatID: "pipit-3", Station: "YST", AOS: night(9, 2, 50), LOS: night(9, 11, 38)},
		{ID: "bk-5002", SatID: "furrow-a", Station: "SQT", AOS: night(9, 40, 30), LOS: night(9, 50, 0)},
		{ID: "bk-5003", SatID: "brine-3", Station: "PEL", AOS: night(10, 15, 0), LOS: night(10, 23, 45)},
	}
}

// readings returns each station's weather at 03:50 on the night.
func readings() map[string]Weather {
	at := night(3, 50, 0)
	return map[string]Weather{
		"YST": {Station: "YST", At: at, WindKmh: 41, GustKmh: 55, TempC: -14},
		"PEL": {Station: "PEL", At: at, WindKmh: 22, GustKmh: 30, TempC: 2},
		"SQT": {Station: "SQT", At: at, WindKmh: 12, GustKmh: 18, TempC: 6},
		"LAN": {Station: "LAN", At: at, WindKmh: 18, GustKmh: 25, TempC: 27},
		"CWR": {Station: "CWR", At: at, WindKmh: 46, GustKmh: 58, TempC: 7},
	}
}

// world is the Board's surroundings in a test: its clock, passplan and the weather service.
type world struct {
	mu      sync.Mutex
	now     time.Time
	passes  []Booking          // what passplan sends: the same slice every time
	weather map[string]Weather // by station code; a station missing here has no reading
}

// newWorld returns the night at the moment at, with its bookings and the 03:50 weather.
func newWorld(at time.Time) *world {
	return &world{now: at, passes: tonight(), weather: readings()}
}

// Now is the Board's clock.
func (w *world) Now() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.now
}

// Passes is passplan's schedule.
func (w *world) Passes(ctx context.Context) ([]Booking, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.passes, nil
}

// Weather is the weather service: a station without a reading answers as a silent mast does.
func (w *world) Weather(ctx context.Context, code string) (Weather, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	r, ok := w.weather[code]
	if !ok {
		return Weather{}, fmt.Errorf("weather for %s: 503 Service Unavailable: mast silent", code)
	}
	return r, nil
}

// board returns the Board in w, with the station table and the rota, refreshed once.
func board(t *testing.T, w *world) *Board {
	t.Helper()
	b, err := NewBoard(Sources{Passes: w.Passes, Weather: w.Weather, Rota: rota()}, stationTable(), w.Now)
	if err != nil {
		t.Fatalf("NewBoard = %v, want no error", err)
	}
	if err := b.Refresh(t.Context()); err != nil {
		t.Fatalf("Refresh = %v, want nil", err)
	}
	return b
}

// get sends GET path to the Board's handler and returns the answer.
func get(b *Board, path string) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	b.Handler().ServeHTTP(res, httptest.NewRequest(http.MethodGet, path, nil))
	return res
}

// lines returns the lines of page that start with prefix, trimmed of spaces and tabs.
func lines(page, prefix string) []string {
	var found []string
	for line := range strings.Lines(page) {
		if line = strings.TrimSpace(line); strings.HasPrefix(line, prefix) {
			found = append(found, line)
		}
	}
	return found
}

// first returns the first line of page that starts with prefix, or "" when none does.
func first(page, prefix string) string {
	if found := lines(page, prefix); len(found) > 0 {
		return found[0]
	}
	return ""
}

// cells returns the text in each cell of a pass row.
func cells(row string) []string {
	row = strings.TrimPrefix(row, `<tr class="pass"><td>`)
	row = strings.TrimSuffix(row, "</td></tr>")
	return strings.Split(row, "</td><td>")
}

// ids returns the booking ID in each pass row of page, in order.
func ids(page string) []string {
	var found []string
	for _, row := range lines(page, `<tr class="pass">`) {
		if c := cells(row); len(c) == 6 {
			found = append(found, c[4])
		} else {
			found = append(found, row)
		}
	}
	return found
}

// row returns the pass row of the booking with this ID on page, or "" when there is none.
func row(page, id string) string {
	for _, r := range lines(page, `<tr class="pass">`) {
		if c := cells(r); len(c) == 6 && c[4] == id {
			return r
		}
	}
	return ""
}

// bookingIDs returns the IDs of bs, in order.
func bookingIDs(bs []Booking) []string {
	var found []string
	for _, b := range bs {
		found = append(found, b.ID)
	}
	return found
}

func TestPage(t *testing.T) {
	res := get(board(t, newWorld(night(3, 10, 0))), "/")
	page := res.Body.String()
	if res.Code != http.StatusOK {
		t.Fatalf("GET / = %d, want 200", res.Code)
	}
	if got, want := res.Header().Get("Content-Type"), "text/html; charset=utf-8"; got != want {
		t.Errorf("GET / Content-Type = %q, want %q", got, want)
	}
	if got, want := first(page, "<title>"), "<title>The Board</title>"; got != want {
		t.Errorf("GET / title line = %q, want %q", got, want)
	}
	if got, want := first(page, "<h1>"), "<h1>The Board · 03:10 UTC</h1>"; got != want {
		t.Errorf("GET / at 03:10 heading = %q, want %q", got, want)
	}
	want := []string{"bk-4986", "bk-4988", "bk-4990", "bk-4991", "bk-4992", "bk-4993", "bk-4994"}
	if got := ids(page); !slices.Equal(got, want) {
		t.Errorf("GET / pass rows = %q, want %q", got, want)
	}
	pipit := `<tr class="pass"><td>04:12</td><td>04:21</td><td>PEL</td><td>pipit-1</td><td>bk-4988</td><td>echo test</td></tr>`
	if got := row(page, "bk-4988"); got != pipit {
		t.Errorf("GET / row of bk-4988 = %q, want %q", got, pipit)
	}
}

func TestPageError(t *testing.T) {
	b := board(t, newWorld(night(3, 10, 0)))
	// This template writes a title and a line per pass, then fails: the view has 7 passes, not 10.
	b.page = template.Must(template.New("broken").Parse(
		"<title>The Board</title>\n{{range .Passes}}<p>{{.SatID}}</p>\n{{end}}{{(index .Passes 9).ID}}\n"))
	res := get(b, "/")
	if res.Code != http.StatusInternalServerError {
		t.Errorf("GET / with a template that fails = %d, want 500", res.Code)
	}
	body := res.Body.String()
	switch {
	case strings.Contains(body, "The Board") || strings.Contains(body, "pipit-1"):
		t.Errorf("GET / with a template that fails sent %q, want none of the page", body)
	case strings.TrimSpace(body) != "internal error":
		t.Errorf("GET / with a template that fails sent %q, want %q", body, "internal error\n")
	}
}

func TestNextPasses(t *testing.T) {
	// shuffled returns the night's bookings in the order passplan happened to send them.
	shuffled := func() []Booking {
		all := tonight()
		return []Booking{all[4], all[0], all[6], all[2], all[1], all[5], all[3]}
	}
	tie := Booking{ID: "bk-4995", SatID: "pipit-3", Station: "LAN", AOS: night(4, 31, 10), LOS: night(4, 38, 40)}
	fifteen := morning()
	slices.Reverse(fifteen)
	fifteen = append(fifteen, shuffled()...)
	for _, tc := range []struct {
		name   string
		at     time.Time
		passes []Booking
		want   []string
	}{
		{"in AOS order", night(3, 10, 0), shuffled(),
			[]string{"bk-4986", "bk-4988", "bk-4990", "bk-4991", "bk-4992", "bk-4993", "bk-4994"}},
		{"the same AOS, by booking ID", night(3, 10, 0), append([]Booking{tie}, shuffled()...),
			[]string{"bk-4986", "bk-4988", "bk-4990", "bk-4995", "bk-4991", "bk-4992", "bk-4993", "bk-4994"}},
		{"one ending now, one under way", night(5, 22, 10), shuffled(),
			[]string{"bk-4993", "bk-4994"}},
		{"twelve at most", night(3, 10, 0), fifteen,
			[]string{"bk-4986", "bk-4988", "bk-4990", "bk-4991", "bk-4992", "bk-4993", "bk-4994", "bk-4996", "bk-4997", "bk-4998", "bk-4999", "bk-5000"}},
		{"the night is over", night(6, 0, 0), shuffled(), nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := newWorld(tc.at)
			w.passes = tc.passes
			sent := bookingIDs(tc.passes)
			page := get(board(t, w), "/").Body.String()
			if got := ids(page); !slices.Equal(got, tc.want) {
				t.Errorf("GET / at %s: pass rows = %q, want %q", tc.at.Format("15:04:05"), got, tc.want)
			}
			if got := bookingIDs(w.passes); !slices.Equal(got, sent) {
				t.Errorf("after Refresh, passplan's slice = %q, want %q as it was sent", got, sent)
			}
		})
	}
}
