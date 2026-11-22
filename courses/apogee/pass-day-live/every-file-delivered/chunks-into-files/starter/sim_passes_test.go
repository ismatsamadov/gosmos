package main

// The simulated world's passes: Season 7's pass model (model.go, each top-level name with
// the sim prefix so it can't collide with the rehearsal's), the pass archive passplan
// writes from it, the plan Season 7's rule makes of that archive, and the week's bookings.

import (
	"bytes"
	"cmp"
	"fmt"
	"math"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// The pass model Season 3's reviewers settled: a circular orbit around a round Earth, the
// orbit's plane turning with the mean Sun for a sun-synchronous satellite (one tilted past
// 90°) and by Earth's bulge for any other, and Earth turning beneath it. A pass is the time the
// satellite's elevation over a station stays at or above the station's mask: its AOS and LOS
// are the instants the elevation crosses the mask, rounded to the second, and its maximum
// elevation is the highest of its elevations a second apart from AOS, to a tenth of a degree.

const (
	simEarthRadiusKm = 6371.0
	simMuKm3PerS2    = 398600.4418
	simJ2            = 1.08263e-3                       // how strongly Earth's bulge turns an orbit
	simEarthRate     = 2 * math.Pi / 86164.0905         // Earth's turn, radians a second
	simSunRate       = 2 * math.Pi / (365.2422 * 86400) // the mean Sun's, radians a second
)

func simRad(d float64) float64 { return d * math.Pi / 180 }
func simDeg(r float64) float64 { return r * 180 / math.Pi }

// simPeriodS returns an orbit's period in seconds, as its altitude gives it.
func simPeriodS(o Orbit) float64 {
	r := simEarthRadiusKm + o.AltKm
	return 2 * math.Pi * math.Sqrt(r*r*r/simMuKm3PerS2)
}

// simSkyTrack is one satellite over one station: everything simPredict needs, worked out once.
// Its times are seconds after the satellite's Epoch.
type simSkyTrack struct {
	a, cosI, sinI   float64    // the orbit's radius in km, and its tilt
	n, u0           float64    // radians a second along the orbit, and how far along at Epoch
	node0, nodeRate float64    // the node's longitude at Epoch, and how fast it moves over Earth
	st, up          [3]float64 // the station, km from Earth's center, and straight up there
	mask, reach     float64    // the mask, and the widest angle at Earth's center it sees out to
	speed           float64    // no satellite and station draw apart or together faster, radians a second
}

func simNewSkyTrack(e Elements, s Site) *simSkyTrack {
	k := &simSkyTrack{a: simEarthRadiusKm + e.Orbit.AltKm, mask: s.MaskDeg}
	i := simRad(e.Orbit.IncDeg)
	k.cosI, k.sinI = math.Cos(i), math.Sin(i)
	k.n = math.Sqrt(simMuKm3PerS2 / (k.a * k.a * k.a))
	if e.RevsPerDay > 0 {
		k.n = e.RevsPerDay * 2 * math.Pi / 86400
	}
	drift := -1.5 * k.n * simJ2 * (simEarthRadiusKm / k.a) * (simEarthRadiusKm / k.a) * k.cosI
	if e.Orbit.IncDeg > 90 {
		drift = simSunRate // sun-synchronous: the plane keeps its angle to the Sun
	}
	k.node0, k.nodeRate, k.u0 = simRad(e.NodeDeg), drift-simEarthRate, simRad(e.PhaseDeg)
	la, lo := simRad(s.LatDeg), simRad(s.LonDeg)
	k.up = [3]float64{math.Cos(la) * math.Cos(lo), math.Cos(la) * math.Sin(lo), math.Sin(la)}
	for j := range 3 {
		k.st[j] = simEarthRadiusKm * k.up[j]
	}
	m := simRad(s.MaskDeg)
	k.reach = math.Acos(simEarthRadiusKm*math.Cos(m)/k.a) - m
	k.speed = k.n + simEarthRate + math.Abs(drift)
	return k
}

// pos is the satellite's place over the turning Earth, km from its center.
func (k *simSkyTrack) pos(t float64) (x, y, z float64) {
	su, cu := math.Sincos(k.u0 + k.n*t)
	sn, cn := math.Sincos(k.node0 + k.nodeRate*t)
	return k.a * (cn*cu - sn*su*k.cosI), k.a * (sn*cu + cn*su*k.cosI), k.a * su * k.sinI
}

// elev is the satellite's elevation over the station at t, in degrees.
func (k *simSkyTrack) elev(t float64) float64 {
	x, y, z := k.pos(t)
	dx, dy, dz := x-k.st[0], y-k.st[1], z-k.st[2]
	d := math.Sqrt(dx*dx + dy*dy + dz*dz)
	return simDeg(math.Asin((dx*k.up[0] + dy*k.up[1] + dz*k.up[2]) / d))
}

// apart is the angle at Earth's center between the satellite and the station.
func (k *simSkyTrack) apart(t float64) float64 {
	x, y, z := k.pos(t)
	return math.Acos(max(-1, min(1, (x*k.up[0]+y*k.up[1]+z*k.up[2])/k.a)))
}

// above is the elevation over the mask: negative while the satellite is below it.
func (k *simSkyTrack) above(t float64) float64 { return k.elev(t) - k.mask }

// orbit is the number of the orbit on which the satellite, near t, crosses the station's
// latitude: Rev0 at Epoch, one more at each northbound equator crossing. A satellite whose
// track never reaches the station's latitude counts from the top or bottom of its track.
func (k *simSkyTrack) orbit(rev0 int, t float64) int {
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

// simPredict returns e's passes over s whose AOS is in [from, to), by AOS.
// No pass is before Epoch or before the station goes live, and a pass shorter than two
// minutes doesn't count. Near the top of its track a satellite crosses a high station's
// latitude twice on one orbit; the station sees it rise once and set once: one pass.
func simPredict(e Elements, s Site, from, to time.Time) []PassView {
	from = simLater(from, simLater(e.Epoch, s.Live))
	if !from.Before(to) {
		return nil
	}
	k := simNewSkyTrack(e, s)
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
			peak, top := k.highest(aos, int(2*k.reach/(k.n-simEarthRate))+60)
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
func (k *simSkyTrack) crossing(lo, hi float64) float64 {
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
func (k *simSkyTrack) highest(aos float64, n int) (int, float64) {
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

// simLater returns the later of two times.
func simLater(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

// simEarlier returns the earlier of two times.
func simEarlier(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

// simNet returns the network every test shares: twenty stations and eighty satellites.
var simNet = sync.OnceValue(newNetwork)

// simWeek is the first instant of the investors' week: Monday 23 July 2029, 00:00 UTC.
var simWeek = time.Date(2029, time.July, 23, 0, 0, 0, 0, time.UTC)

// simJuly returns h:m:s on a day of July 2029, UTC: simJuly(26, 16, 3, 46) is Thursday's
// 16:03:46 in the investors' week.
func simJuly(day, h, m, s int) time.Time {
	return time.Date(2029, time.July, day, h, m, s, 0, time.UTC)
}

// simWindow returns the window that starts at from and lasts days days.
func simWindow(from time.Time, days int) Window {
	return Window{From: from, To: from.Add(time.Duration(days) * 24 * time.Hour)}
}

// simByAOS orders passes by AOS, then by pass ID, as Season 7's byAOS does.
func simByAOS(a, b PassView) int {
	return cmp.Or(a.AOS.Compare(b.AOS), strings.Compare(a.ID, b.ID))
}

// simSite returns the station with the given code.
func simSite(code string) (Site, bool) {
	for _, s := range simNet().Sites {
		if s.Code == code {
			return s, true
		}
	}
	return Site{}, false
}

// simSatellite returns the network's satellite with the given ID.
func simSatellite(id string) (Satellite, bool) {
	for _, s := range simNet().Satellites {
		if s.ID == id {
			return s, true
		}
	}
	return Satellite{}, false
}

// simCodes returns the given station codes in table order, or all twenty for none.
func simCodes(stations []string) []string {
	var out []string
	for _, s := range simNet().Sites {
		if len(stations) == 0 || slices.Contains(stations, s.Code) {
			out = append(out, s.Code)
		}
	}
	return out
}

// simCanServe reports whether the station offers one of the satellite's bands.
func simCanServe(s Site, sat Satellite) bool {
	return slices.ContainsFunc(sat.Bands, func(b Band) bool { return slices.Contains(s.Bands, b) })
}

// simPasses caches what the model predicts, for the life of the test binary: predicting
// the investors' week for twenty stations takes a moment, and many cases ask for it.
var simPasses struct {
	mu       sync.Mutex
	stations map[string][]PassView // one station's passes in a window, by AOS then pass ID
	archives map[string][]PassView
	texts    map[string][]byte
	plans    map[string]*simPlan
}

// simKey names a window and a set of stations for the cache.
func simKey(win Window, codes []string) string {
	return win.From.Format(time.RFC3339) + " " + win.To.Format(time.RFC3339) + " " + strings.Join(codes, ",")
}

// simPredictStations predicts every satellite's passes over each of the stations in the
// window, with as many workers as there are CPUs, and keeps each station's list.
func simPredictStations(win Window, codes []string) {
	type job struct {
		site Site
		el   Elements
	}
	n := simNet()
	var jobs []job
	for _, code := range codes {
		if _, ok := simPasses.stations[simKey(win, []string{code})]; ok {
			continue
		}
		site, _ := simSite(code)
		for _, e := range n.Elements {
			jobs = append(jobs, job{site, e})
		}
	}
	if len(jobs) == 0 {
		return
	}
	results := make([][]PassView, len(jobs))
	next := make(chan int)
	var wg sync.WaitGroup
	for range runtime.GOMAXPROCS(0) {
		wg.Go(func() {
			for j := range next {
				results[j] = simPredict(jobs[j].el, jobs[j].site, win.From, win.To)
			}
		})
	}
	for j := range jobs {
		next <- j
	}
	close(next)
	wg.Wait()
	if simPasses.stations == nil {
		simPasses.stations = map[string][]PassView{}
	}
	for j, r := range results {
		key := simKey(win, []string{jobs[j].site.Code})
		simPasses.stations[key] = append(simPasses.stations[key], r...)
	}
	for _, code := range codes {
		slices.SortFunc(simPasses.stations[simKey(win, []string{code})], simByAOS)
	}
}

// simArchive returns passplan's pass archive for the window: every predicted pass of every
// satellite over each of the stations (any band, any elevation from the mask up, two
// minutes or longer) with its AOS in the window, by AOS then pass ID. Nobody changes it.
func simArchive(win Window, stations []string) []PassView {
	codes := simCodes(stations)
	key := simKey(win, codes)
	simPasses.mu.Lock()
	defer simPasses.mu.Unlock()
	if a, ok := simPasses.archives[key]; ok {
		return a
	}
	simPredictStations(win, codes)
	var all []PassView
	for _, code := range codes {
		all = append(all, simPasses.stations[simKey(win, []string{code})]...)
	}
	slices.SortFunc(all, simByAOS)
	if simPasses.archives == nil {
		simPasses.archives = map[string][]PassView{}
	}
	simPasses.archives[key] = all
	return all
}

// simArchiveLine writes one pass as passes.txt holds it, without the newline.
func simArchiveLine(p PassView) string {
	return p.ID + " " + p.AOS.Format(time.RFC3339) + " " + p.LOS.Format(time.RFC3339) + " " +
		strconv.FormatFloat(p.MaxElevDeg, 'f', 1, 64)
}

// simArchiveText returns the archive for the window and stations as the bytes of
// passes.txt. Nobody changes it.
func simArchiveText(win Window, stations []string) []byte {
	passes := simArchive(win, stations)
	key := simKey(win, simCodes(stations))
	simPasses.mu.Lock()
	defer simPasses.mu.Unlock()
	if t, ok := simPasses.texts[key]; ok {
		return t
	}
	var buf bytes.Buffer
	buf.Grow(70 * len(passes))
	for _, p := range passes {
		buf.WriteString(simArchiveLine(p))
		buf.WriteByte('\n')
	}
	if simPasses.texts == nil {
		simPasses.texts = map[string][]byte{}
	}
	simPasses.texts[key] = buf.Bytes()
	return buf.Bytes()
}

// simPlan is what Season 7's rule makes of an archive: the reference every rehearsal's
// plan is compared with.
type simPlan struct {
	win        Window
	contacts   []Contact          // by AOS, then pass ID
	byID       map[string]Contact // the same contacts by pass ID
	candidates int                // servable passes of at least 10 degrees
	dropped    int                // candidates no dish could take
}

// simPlanFor plans the window at the stations by Season 7's rule: each station on its own,
// its candidates in AOS order then pass ID, each taking the lowest-numbered dish free at
// its AOS, a dish being free two minutes after its last contact's LOS. Nobody changes it.
func simPlanFor(win Window, stations []string) *simPlan {
	archive := simArchive(win, stations)
	codes := simCodes(stations)
	key := simKey(win, codes)
	simPasses.mu.Lock()
	defer simPasses.mu.Unlock()
	if p, ok := simPasses.plans[key]; ok {
		return p
	}
	plan := &simPlan{win: win, byID: map[string]Contact{}}
	perStation := map[string][]PassView{}
	for _, p := range archive {
		site, _ := simSite(p.Station)
		sat, _ := simSatellite(p.SatID)
		if simCanServe(site, sat) && p.MaxElevDeg >= 10 {
			perStation[p.Station] = append(perStation[p.Station], p)
		}
	}
	for _, code := range codes {
		site, _ := simSite(code)
		free := make([]time.Time, site.Dishes)
		for _, p := range perStation[code] {
			plan.candidates++
			d := slices.IndexFunc(free, func(t time.Time) bool { return !p.AOS.Before(t) })
			if d < 0 {
				plan.dropped++
				continue
			}
			free[d] = p.LOS.Add(2 * time.Minute)
			plan.contacts = append(plan.contacts, Contact{PassView: p, Dish: d + 1})
		}
	}
	slices.SortFunc(plan.contacts, func(a, b Contact) int { return simByAOS(a.PassView, b.PassView) })
	for _, c := range plan.contacts {
		plan.byID[c.ID] = c
	}
	if simPasses.plans == nil {
		simPasses.plans = map[string]*simPlan{}
	}
	simPasses.plans[key] = plan
	return plan
}

// days returns the window's days as plan files name them: "2029-07-23", ...
func (p *simPlan) days() []string {
	var out []string
	for d := p.win.From; d.Before(p.win.To); d = d.Add(24 * time.Hour) {
		out = append(out, d.Format("2006-01-02"))
	}
	return out
}

// dayLines returns the lines plan/<day>.txt must hold: the contacts whose AOS falls on that
// UTC day, by AOS then pass ID, as "16:03:46 16:09:59 1 thistle-2/PEL/02136".
func (p *simPlan) dayLines(day string) []string {
	var out []string
	for _, c := range p.contacts {
		if c.AOS.Format("2006-01-02") == day {
			out = append(out, fmt.Sprintf("%s %s %d %s", c.AOS.Format("15:04:05"), c.LOS.Format("15:04:05"), c.Dish, c.ID))
		}
	}
	return out
}

// simFirstScience is the AOS of Thistle-2's first science pass, thistle-2/PEL/02136:
// Kessick Polytechnic books every Thistle-2 contact from it on.
var simFirstScience = simJuly(26, 16, 3, 46)

// simBook returns the bookings customers hold on a plan, numbered from bk-8001 in the
// plan's order: every contact of a Furrowlight, Brinewatch or Cirrine Atmospherics
// satellite with a maximum elevation of at least 70 degrees, and every Thistle-2 contact
// from its first science pass on. Apogee books none of its own passes. For the investors'
// week at all twenty stations these are bk-8001 to bk-9493.
func simBook(contacts []Contact) []Booking {
	var out []Booking
	for _, c := range contacts {
		sat, _ := simSatellite(c.SatID)
		booked := false
		switch sat.Operator {
		case "Apogee":
		case "Kessick Polytechnic":
			booked = !c.AOS.Before(simFirstScience)
		default:
			booked = c.MaxElevDeg >= 70
		}
		if booked {
			id := fmt.Sprintf("bk-%04d", 8001+len(out))
			out = append(out, Booking{ID: id, SatID: c.SatID, Station: c.Station, AOS: c.AOS, LOS: c.LOS})
		}
	}
	return out
}
