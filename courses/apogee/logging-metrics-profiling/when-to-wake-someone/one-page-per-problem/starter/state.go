package main

import "time"

// ruleState is where one rule stands at one station.
type ruleState struct {
	since  time.Time // when the value reached FireAt; zero while it is below
	firing bool      // fired, and not yet below ClearBelow
	sent   time.Time // when FIRING was last sent
}

// beforeFiring takes a sample while the rule isn't firing, and returns "FIRING" when the rule
// fires, or "" for no alert.
func (st *ruleState) beforeFiring(r Rule, s Sample) string {
	if s.Value < r.FireAt {
		st.since = time.Time{}
		return ""
	}
	if st.since.IsZero() {
		st.since = s.At
	}
	if s.At.Sub(st.since) < r.For {
		return ""
	}
	st.firing, st.sent = true, s.At
	return "FIRING"
}

// whileFiring takes a sample while the rule is firing, and returns "FIRING", "RESOLVED" or ""
// for no alert.
func (st *ruleState) whileFiring(r Rule, s Sample) string {
	if s.Value < r.ClearBelow {
		*st = ruleState{}
		return "RESOLVED"
	}
	return "FIRING"
}
