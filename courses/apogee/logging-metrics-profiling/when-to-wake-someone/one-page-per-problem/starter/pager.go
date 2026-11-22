package main

// Pager keeps each rule's state per station.
type Pager struct {
	Rules    []Rule
	Rota     []Shift
	Silences []Silence
	state    map[[2]string]*ruleState // by rule name and station
}

// Observe takes the next sample of rule's metric at one station and returns the alert to
// send, if any.
func (p *Pager) Observe(rule string, s Sample) (Alert, bool) {
	r, ok := p.rule(rule)
	if !ok {
		return Alert{}, false
	}
	st := p.stateOf(rule, s.Station)
	var kind string
	if st.firing {
		kind = st.whileFiring(r, s)
	} else {
		kind = st.beforeFiring(r, s)
	}
	if kind == "" {
		return Alert{}, false
	}
	return Alert{At: s.At, Rule: rule, Station: s.Station, Kind: kind, Value: s.Value}, true
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
