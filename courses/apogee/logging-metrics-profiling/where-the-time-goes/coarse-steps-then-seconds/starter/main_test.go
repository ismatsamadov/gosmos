package main

import (
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"
)

// from and to are the twelve hours the tests plan: Tuesday 16 January from 05:10, when Cape Wren stowed.
var (
	from = time.Date(2029, 1, 16, 5, 10, 0, 0, time.UTC)
	to   = from.Add(12 * time.Hour)
)

// show writes contacts as their AOS and LOS to the second, or says there are none.
func show(cs []contact) string {
	if len(cs) == 0 {
		return "none"
	}
	var s []string
	for _, c := range cs {
		s = append(s, toSecond(c.AOS).Format("15:04:05")+"-"+toSecond(c.LOS).Format("15:04:05"))
	}
	return strings.Join(s, " ")
}

// counted returns sat's elevation over st as a function, and the number of times it was asked.
func counted(sat Satellite, st Station) (func(time.Time) float64, *int) {
	n := 0
	elems := elements()[sat.ID]
	return func(at time.Time) float64 { n++; return elevationDeg(sat, elems, st, at) }, &n
}

func TestFindPasses(t *testing.T) {
	for _, sat := range fleet() {
		for _, code := range []string{"CWR", "LAN"} {
			st := station(code)
			ok := t.Run(sat.ID+" over "+code, func(t *testing.T) {
				scan, asked := counted(sat, st)
				want := show(scanEverySecond(scan, from, to, st.MaskDeg))
				el, calls := counted(sat, st)
				got := show(findPasses(el, from, to, st.MaskDeg))
				call := fmt.Sprintf("findPasses(%s over %s, 05:10, 17:10, %v)", sat.ID, code, st.MaskDeg)
				if got != want {
					t.Errorf("%s = %s, want %s", call, got, want)
				}
				if *calls*20 > *asked {
					t.Errorf("%s asked el %d times, want at most %d, a twentieth of the scan's %d", call, *calls, *asked/20, *asked)
				}
			})
			if !ok {
				return
			}
		}
	}
}

func TestFindPassesEdges(t *testing.T) {
	sat, st := satellite("brine-3"), station("CWR")
	at := func(hms string) time.Time {
		v, _ := time.Parse(time.RFC3339, "2029-01-16T"+hms+"Z")
		return v
	}
	tests := []struct {
		name       string
		start, end time.Time
	}{
		{"from inside the 05:29 pass to inside the 07:07 pass", at("05:33:00"), at("07:12:00")},
		{"ending four seconds after the 07:07 pass's LOS", at("05:10:00"), at("07:16:55")},
		{"ending at 05:39:00, with the 05:29 pass still up at 05:38:59", at("05:10:00"), at("05:39:00")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scan, _ := counted(sat, st)
			el, _ := counted(sat, st)
			got, want := show(findPasses(el, tt.start, tt.end, st.MaskDeg)), show(scanEverySecond(scan, tt.start, tt.end, st.MaskDeg))
			if got != want {
				t.Errorf("findPasses(brine-3 over CWR, %s, %s, 5) = %s, want %s", tt.start.Format("15:04:05"), tt.end.Format("15:04:05"), got, want)
			}
		})
	}
}

func TestFindPassesSpeed(t *testing.T) {
	flag.Set("test.benchtime", "100ms")
	sat, st := satellite("brine-3"), station("CWR")
	el, _ := counted(sat, st)
	scan := testing.Benchmark(func(b *testing.B) {
		for b.Loop() {
			scanEverySecond(el, from, to, st.MaskDeg)
		}
	})
	coarse := testing.Benchmark(func(b *testing.B) {
		for b.Loop() {
			findPasses(el, from, to, st.MaskDeg)
		}
	})
	if coarse.NsPerOp()*5 > scan.NsPerOp() {
		t.Errorf("findPasses(brine-3 over CWR, 05:10, 17:10, 5) took %d ns a run, want at most %d, a fifth of scanEverySecond's %d", coarse.NsPerOp(), scan.NsPerOp()/5, scan.NsPerOp())
	}
}
