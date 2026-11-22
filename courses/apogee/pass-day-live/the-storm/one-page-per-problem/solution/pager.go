package main

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"sort"
	"time"
)

// Severity ranks an alert.
type Severity int

const (
	Info Severity = iota
	Warning
	Critical
)

// Rule is one of the pager's pass rules: a station's failed passes over a span, watched
// every minute.
type Rule struct {
	Name       string        // "pass-failures"
	Over       time.Duration // the failures counted: those with LOS in (now − Over, now]
	FireAt     float64       // the problem starts when the count is at or above this
	ClearBelow float64       // and is over only when it falls below this
	Repeat     time.Duration // a problem still firing is sent again after this long; 0: never
	Severity   Severity      // Critical pages whoever is on call; Warning goes on the morning list
}

// passRules are the pager's pass rules since 8 January 2029, by name.
func passRules() []Rule {
	return []Rule{
		{Name: "failed-pass", Over: time.Minute, FireAt: 1, ClearBelow: 1, Severity: Warning},
		{Name: "pass-failures", Over: 6 * time.Hour, FireAt: 2, ClearBelow: 1, Repeat: 24 * time.Hour, Severity: Critical},
	}
}

// Alert is one message the pager sends.
type Alert struct {
	At      time.Time
	Rule    string
	Station string
	Kind    string // "FIRING" or "RESOLVED"
	Value   float64
	To      string // the person on call, or "morning list"
}

// String is the alert as alerts.txt writes it: "Thu 10:56 FIRING pass-failures YST 2 to No. 039".
func (a Alert) String() string {
	return fmt.Sprintf("%s %s %s %s %g to %s", a.At.Format("Mon 15:04"), a.Kind, a.Rule, a.Station, a.Value, a.To)
}

// Pager watches every rehearsed station's failed passes and sends alerts by the rules.
type Pager struct {
	r       *run
	rules   []Rule
	rota    []Shift
	los     []time.Time            // every contact's LOS, ascending
	results chan passResult        // each contact's, at its LOS
	failed  map[string][]time.Time // each station's failed passes, by LOS
	firing  map[[2]string]firing   // by rule and station: the problems firing now
	alerts  []Alert
}

// firing is a problem the pager has sent: when its last FIRING went out, and to whom.
type firing struct {
	at time.Time
	to string
}

// passResult is how one contact ended, as its station tells the pager at LOS.
type passResult struct {
	station string
	los     time.Time
	failed  bool
}

// newPager returns the pager for a week.
func newPager(r *run) *Pager {
	p := &Pager{r: r, rules: passRules(), rota: onCall(), results: make(chan passResult, len(r.wk.Contacts)),
		failed: map[string][]time.Time{}, firing: map[[2]string]firing{}}
	for _, c := range r.wk.Contacts {
		p.los = append(p.los, c.LOS)
	}
	slices.SortFunc(p.los, time.Time.Compare)
	return p
}

// result tells the pager how a contact ended. Skipped and stopped passes are not failures.
func (p *Pager) result(c Contact, o Outcome) {
	p.results <- passResult{c.Station, c.LOS, o.Kind == "failed"}
}

// run evaluates the rules on every minute of the clock until the week is done; the instant
// it is done comes on final, and the minutes up to it are evaluated first.
func (p *Pager) run(ctx context.Context, final <-chan time.Time) {
	w := p.r.w
	next := w.Now().Truncate(time.Minute).Add(time.Minute)
	got := 0
	for {
		t := time.NewTimer(next.Sub(w.Now()))
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case end := <-final:
			t.Stop()
			for ; !next.After(end); next = next.Add(time.Minute) {
				p.minute(ctx, next, &got)
			}
			return
		case <-t.C:
			if !p.minute(ctx, next, &got) {
				return
			}
			next = next.Add(time.Minute)
		}
	}
}

// minute evaluates every rule for every station at one minute, once every contact whose LOS
// is not after it has told its result. It reports false when ctx ended first.
func (p *Pager) minute(ctx context.Context, at time.Time, got *int) bool {
	due := sort.Search(len(p.los), func(i int) bool { return p.los[i].After(at) })
	for *got < due {
		select {
		case res := <-p.results:
			*got++
			if res.failed {
				p.failed[res.station] = append(p.failed[res.station], res.los)
			}
		case <-ctx.Done():
			return false
		}
	}
	for _, s := range p.r.wk.Sites {
		for _, rule := range p.rules {
			p.evaluate(rule, s.Code, at)
		}
	}
	return true
}

// evaluate moves one rule's state at one station on one minute, and sends what it says.
func (p *Pager) evaluate(rule Rule, station string, at time.Time) {
	value := 0.0
	for _, los := range p.failed[station] {
		if los.After(at.Add(-rule.Over)) && !los.After(at) {
			value++
		}
	}
	key := [2]string{rule.Name, station}
	f, on := p.firing[key]
	switch {
	case !on && value >= rule.FireAt, on && rule.Repeat > 0 && at.Sub(f.at) >= rule.Repeat:
		p.firing[key] = firing{at, p.send(rule, station, "FIRING", value, at, "")}
	case on && value < rule.ClearBelow:
		delete(p.firing, key)
		if rule.Severity == Critical {
			p.send(rule, station, "RESOLVED", value, at, f.to)
		}
	}
}

// send sends an alert and returns its recipient: a Critical FIRING goes to whoever is on
// call at its minute, a RESOLVED to the person its FIRING went to, and a Warning to the
// morning list.
func (p *Pager) send(rule Rule, station, kind string, value float64, at time.Time, to string) string {
	switch {
	case rule.Severity != Critical:
		to = "morning list"
	case to == "":
		for _, s := range p.rota {
			if !s.From.After(at) {
				to = s.Name
			}
		}
	}
	a := Alert{At: at, Rule: rule.Name, Station: station, Kind: kind, Value: value, To: to}
	p.alerts = append(p.alerts, a)
	p.r.logs.alert(a)
	return to
}

// write saves reports/alerts.txt in dir: one line per alert, in time order.
func (p *Pager) write(dir string) error {
	var buf bytes.Buffer
	for _, a := range p.alerts {
		fmt.Fprintln(&buf, a)
	}
	return saveWhole(filepath.Join(dir, "reports", "alerts.txt"), buf.Bytes())
}
