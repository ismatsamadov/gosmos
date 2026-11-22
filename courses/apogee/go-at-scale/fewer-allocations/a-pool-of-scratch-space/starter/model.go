package main

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// The pass model Season 3's reviewers settled: a circular orbit around a round Earth, the
// orbit's plane turning with the mean Sun for a sun-synchronous satellite (one tilted past
// 90°) and by Earth's bulge for any other, and Earth turning beneath it. A pass is the time the
// satellite's elevation over a station stays at or above the station's mask: its AOS and LOS
// are the instants the elevation crosses the mask, rounded to the second, and its maximum
// elevation is the highest of its elevations a second apart from AOS, to a tenth of a degree.

const (
	earthRadiusKm = 6371.0
	muKm3PerS2    = 398600.4418
	j2            = 1.08263e-3                       // how strongly Earth's bulge turns an orbit
	earthRate     = 2 * math.Pi / 86164.0905         // Earth's turn, radians a second
	sunRate       = 2 * math.Pi / (365.2422 * 86400) // the mean Sun's, radians a second
)

func rad(d float64) float64 { return d * math.Pi / 180 }
func deg(r float64) float64 { return r * 180 / math.Pi }

// periodS returns an orbit's period in seconds, as its altitude gives it.
func periodS(o Orbit) float64 {
	r := earthRadiusKm + o.AltKm
	return 2 * math.Pi * math.Sqrt(r*r*r/muKm3PerS2)
}

// skyTrack is one satellite over one station: everything predict needs, worked out once.
// Its times are seconds after the satellite's Epoch.
type skyTrack struct {
	a, cosI, sinI   float64    // the orbit's radius in km, and its tilt
	n, u0           float64    // radians a second along the orbit, and how far along at Epoch
	node0, nodeRate float64    // the node's longitude at Epoch, and how fast it moves over Earth
	st, up          [3]float64 // the station, km from Earth's center, and straight up there
	mask, reach     float64    // the mask, and the widest angle at Earth's center it sees out to
	speed           float64    // no satellite and station draw apart or together faster, radians a second
}

func newSkyTrack(e Elements, s Site) *skyTrack {
	k := &skyTrack{a: earthRadiusKm + e.Orbit.AltKm, mask: s.MaskDeg}
	i := rad(e.Orbit.IncDeg)
	k.cosI, k.sinI = math.Cos(i), math.Sin(i)
	k.n = math.Sqrt(muKm3PerS2 / (k.a * k.a * k.a))
	if e.RevsPerDay > 0 {
		k.n = e.RevsPerDay * 2 * math.Pi / 86400
	}
	drift := -1.5 * k.n * j2 * (earthRadiusKm / k.a) * (earthRadiusKm / k.a) * k.cosI
	if e.Orbit.IncDeg > 90 {
		drift = sunRate // sun-synchronous: the plane keeps its angle to the Sun
	}
	k.node0, k.nodeRate, k.u0 = rad(e.NodeDeg), drift-earthRate, rad(e.PhaseDeg)
	la, lo := rad(s.LatDeg), rad(s.LonDeg)
	k.up = [3]float64{math.Cos(la) * math.Cos(lo), math.Cos(la) * math.Sin(lo), math.Sin(la)}
	for j := range 3 {
		k.st[j] = earthRadiusKm * k.up[j]
	}
	m := rad(s.MaskDeg)
	k.reach = math.Acos(earthRadiusKm*math.Cos(m)/k.a) - m
	k.speed = k.n + earthRate + math.Abs(drift)
	return k
}

// pos is the satellite's place over the turning Earth, km from its center.
func (k *skyTrack) pos(t float64) (x, y, z float64) {
	su, cu := math.Sincos(k.u0 + k.n*t)
	sn, cn := math.Sincos(k.node0 + k.nodeRate*t)
	return k.a * (cn*cu - sn*su*k.cosI), k.a * (sn*cu + cn*su*k.cosI), k.a * su * k.sinI
}

// elev is the satellite's elevation over the station at t, in degrees.
func (k *skyTrack) elev(t float64) float64 {
	x, y, z := k.pos(t)
	dx, dy, dz := x-k.st[0], y-k.st[1], z-k.st[2]
	d := math.Sqrt(dx*dx + dy*dy + dz*dz)
	return deg(math.Asin((dx*k.up[0] + dy*k.up[1] + dz*k.up[2]) / d))
}

// apart is the angle at Earth's center between the satellite and the station.
func (k *skyTrack) apart(t float64) float64 {
	x, y, z := k.pos(t)
	return math.Acos(max(-1, min(1, (x*k.up[0]+y*k.up[1]+z*k.up[2])/k.a)))
}

// above is the elevation over the mask: negative while the satellite is below it.
func (k *skyTrack) above(t float64) float64 { return k.elev(t) - k.mask }

// orbit is the number of the orbit on which the satellite, near t, crosses the station's
// latitude: Rev0 at Epoch, one more at each northbound equator crossing. A satellite whose
// track never reaches the station's latitude counts from the top or bottom of its track.
func (k *skyTrack) orbit(rev0 int, t float64) int {
	u := k.u0 + k.n*t
	var marks []float64 // where along an orbit the satellite is at the station's latitude
	if s := k.up[2] / k.sinI; s >= 1 {
		marks = []float64{math.Pi / 2}
	} else if s <= -1 {
		marks = []float64{3 * math.Pi / 2}
	} else {
		marks = []float64{math.Asin(s), math.Pi - math.Asin(s)}
	}
	best := math.Inf(1)
	var at float64
	for _, m := range marks {
		c := m + 2*math.Pi*math.Round((u-m)/(2*math.Pi)) // the nearest such point to u
		if d := math.Abs(c - u); d < best {
			best, at = d, c
		}
	}
	return rev0 + int(math.Floor(at/(2*math.Pi)))
}

// predict returns e's passes over s whose AOS is in [from, to), by AOS.
// No pass is before Epoch or before the station goes live, and a pass shorter than two
// minutes doesn't count. Near the top of its track a satellite crosses a high station's
// latitude twice on one orbit; the station sees it rise once and set once: one pass.
func predict(e Elements, s Site, from, to time.Time) []PassView {
	from = later(from, later(e.Epoch, s.Live))
	if !from.Before(to) {
		return nil
	}
	k := newSkyTrack(e, s)
	start, end := from.Sub(e.Epoch).Seconds(), to.Sub(e.Epoch).Seconds()
	at := func(t float64) time.Time { // to the nearest second
		return e.Epoch.Add(time.Duration(math.RoundToEven(t)) * time.Second)
	}
	tenths := func(v float64) float64 { // to a tenth, as the settled model prints it
		f, _ := strconv.ParseFloat(strconv.FormatFloat(v, 'f', 1, 64), 64)
		return f
	}
	var out []PassView
	const step = 20.0
	t := start - 1800 // half an hour early: a pass in progress at from is seen rising
	prev := k.above(t)
	for t < end+60 {
		if prev < 0 {
			// Far below the mask the satellite can't come up for a while: skip what it can't.
			if gap := k.apart(t) - k.reach - 0.01; gap > 0 {
				if jump := math.Floor(gap/k.speed/step) - 1; jump >= 1 {
					t += jump * step
				}
			}
		}
		t += step
		cur := k.above(t)
		if prev < 0 && cur >= 0 {
			aos := k.crossing(t-step, t)
			peak, top := k.highest(aos, int(2*k.reach/(k.n-earthRate))+60)
			end := aos + float64(2*peak) + 120
			for k.above(end) >= 0 {
				end += 60
			}
			los := k.crossing(aos+float64(peak), end)
			p := Pass{SatID: e.SatID, Station: s.Code, AOS: at(aos), LOS: at(los), MaxElevDeg: tenths(top)}
			id := fmt.Sprintf("%s/%s/%05d", e.SatID, s.Code, k.orbit(e.Rev0, aos+float64(peak)))
			out = append(out, PassView{ID: id, Pass: p})
			t = los
			cur = k.above(t)
		}
		prev = cur
	}
	kept := out[:0]
	for _, p := range out {
		if !p.AOS.Before(from) && p.AOS.Before(to) && p.LOS.Sub(p.AOS) >= 2*time.Minute {
			kept = append(kept, p)
		}
	}
	return kept
}

// crossing returns the instant between lo and hi where the elevation crosses the mask, below
// it at one end and not at the other: the Illinois method, to a hundred-millionth of a second.
func (k *skyTrack) crossing(lo, hi float64) float64 {
	f := k.above
	flo, fhi := f(lo), f(hi)
	side := 0
	for range 100 {
		m := (lo*fhi - hi*flo) / (fhi - flo)
		if !(m > lo && m < hi) {
			m = (lo + hi) / 2
		}
		fm := f(m)
		if (fm >= 0) == (fhi >= 0) {
			hi, fhi = m, fm
			if side == 1 {
				flo /= 2
			}
			side = 1
		} else {
			lo, flo = m, fm
			if side == -1 {
				fhi /= 2
			}
			side = -1
		}
		if hi-lo < 1e-8 || fm == 0 {
			return m
		}
	}
	return (lo + hi) / 2
}

// highest returns the whole second i in [0, n] after aos when the elevation is largest, and
// that elevation. The elevation rises to one peak, then falls, so a binary search on which
// way it is heading finds the peak.
func (k *skyTrack) highest(aos float64, n int) (int, float64) {
	f := func(i int) float64 { return k.elev(aos + float64(i)) }
	lo, hi := 0, n
	for lo < hi {
		m := (lo + hi) / 2
		if f(m+1) > f(m) {
			lo = m + 1
		} else {
			hi = m
		}
	}
	return lo, f(lo)
}

// later returns the later of two times.
func later(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// earlier returns the earlier of two times.
func earlier(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
