package main

// Pager keeps each rule's state per station.
type Pager struct {
	Rules    []Rule
	Rota     []Shift
	Silences []Silence
	state    map[[2]string]*ruleState // by rule name and station
}

// Observe takes the next sample of rule's metric at one station and returns the alert to
// send, if any, routed to whoever reads it. While a silence covers the station it sends
// nothing; a FIRING it held back goes out at the first sample once the silence has ended,
// unless the problem is over by then.
func (p *Pager) Observe(rule string, s Sample) (Alert, bool) {
	r, ok := p.rule(rule)
	if !ok {
		return Alert{}, false
	}
	st := p.stateOf(rule, s.Station)
	was := *st // where the rule stood before this sample
	var kind string
	if st.firing {
		kind = st.whileFiring(r, s)
	} else {
		kind = st.beforeFiring(r, s)
	}
	if silenced(p.Silences, s.Station, s.At) {
		if kind == "FIRING" && was.firing {
			st.sent = was.sent // a reminder nobody got: it stays due
		} else if kind == "FIRING" {
			st.held = true // a new problem nobody has heard of
		}
		return Alert{}, false
	}
	switch {
	case st.held: // the problem outlasted the silence: send its FIRING now
		st.held, st.sent, kind = false, s.At, "FIRING"
	case was.held: // it is over, and nobody heard it start: nobody is told it ended
		kind = ""
	}
	if kind == "" {
		return Alert{}, false
	}
	a := Alert{At: s.At, Rule: rule, Station: s.Station, Kind: kind, Value: s.Value}
	return route(a, r.Severity, p.Rota), true
}

// rule returns the pager's rule called name.
func (p *Pager) rule(name string) (Rule, bool) {
	for _, r := range p.Rules {
		if r.Name == name {
			return r, true
		}
	}
	return Rule{}, false
}

// stateOf returns where rule stands at station, starting it inactive at first use.
func (p *Pager) stateOf(rule, station string) *ruleState {
	key := [2]string{rule, station}
	if p.state[key] == nil {
		if p.state == nil {
			p.state = make(map[[2]string]*ruleState)
		}
		p.state[key] = &ruleState{}
	}
	return p.state[key]
}
