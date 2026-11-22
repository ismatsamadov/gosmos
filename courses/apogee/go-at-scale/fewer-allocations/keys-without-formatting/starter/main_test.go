package main

import (
	"slices"
	"testing"
	"time"
)

// aprilSecond returns the candidates of Monday 2 April at every
// station, station by station in table order.
func aprilSecond() []PassView {
	apr2 := time.Date(2029, 4, 2, 0, 0, 0, 0, time.UTC)
	var out []PassView
	for _, s := range newNetwork().Sites {
		_, day := stationDay(s.Code, apr2)
		out = append(out, day...)
	}
	return out
}

// risingAt returns the candidates at the stations with codes that
// rise at the given instant, station by station.
func risingAt(at time.Time, codes ...string) []PassView {
	day := at.Truncate(24 * time.Hour)
	var out []PassView
	for _, code := range codes {
		_, cands := stationDay(code, day)
		for _, p := range cands {
			if p.AOS.Equal(at) {
				out = append(out, p)
			}
		}
	}
	return out
}

func TestDedupe(t *testing.T) {
	apr2 := aprilSecond()
	var inTurn []PassView
	for _, p := range apr2 {
		inTurn = append(inTurn, p, p)
	}
	pel := time.Date(2029, 4, 8, 10, 29, 34, 0, time.UTC)
	twoSats := risingAt(pel, "PEL")
	may17 := time.Date(2029, 5, 17, 7, 8, 9, 0, time.UTC)
	twoStations := risingAt(may17, "TIR", "YST")
	tests := []struct {
		name   string
		passes []PassView
		want   []PassView
	}{
		{"2 April's candidates, then all again", append(slices.Clone(apr2), apr2...), apr2},
		{"2 April's candidates, each twice in a row", inTurn, apr2},
		{"Brine-3 and Pipit-3 over Port Ellery at 10:29:34", twoSats, twoSats},
		{"Pipit-3 over Tirrin Rock and Ystmark at 07:08:09", twoStations, twoStations},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dedupe(slices.Clone(tt.passes))
			if len(got) != len(tt.want) {
				t.Fatalf("len(dedupe(%d passes)) = %d, want %d", len(tt.passes), len(got), len(tt.want))
			}
			for i, p := range tt.want {
				if got[i] != p {
					t.Fatalf("dedupe(%d passes)[%d] = %s, want %s: each pass where it first came", len(tt.passes), i, got[i].ID, p.ID)
				}
			}
		})
	}
}

func TestDedupeAllocs(t *testing.T) {
	if raceOn {
		t.Skip("allocation counts differ under -race")
	}
	once := aprilSecond()
	twice := append(slices.Clone(once), once...)
	var a, b []PassView
	allocsOnce := testing.AllocsPerRun(10, func() { a = dedupe(once) })
	allocsTwice := testing.AllocsPerRun(10, func() { b = dedupe(twice) })
	if !slices.Equal(a, once) || !slices.Equal(b, once) {
		return // TestDedupe says what is wrong
	}
	if extra := allocsTwice - allocsOnce; extra > 0 {
		t.Errorf("%d lookups of passes already seen made %.0f allocations, want 0: dedupe made %.0f for 2 April's candidates twice, %.0f for them once", len(once), extra, allocsTwice, allocsOnce)
	}
}
