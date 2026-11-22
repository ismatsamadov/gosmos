package main

import "time"

// night is how Lantern Atoll's passes ended from Wednesday 27 to Thursday 28 December 2028,
// the night the old rules paged Tomás eleven times.
var night = []Result{
	{PassID: "brine-3/LAN/15145", Station: "LAN", LOS: time.Date(2028, 12, 27, 19, 41, 28, 0, time.UTC), Outcome: "failed", Reason: "deadline"},
	{PassID: "pipit-1/LAN/14031", Station: "LAN", LOS: time.Date(2028, 12, 27, 20, 20, 49, 0, time.UTC), Outcome: "failed", Reason: "deadline"},
	{PassID: "cirrine-1/LAN/17147", Station: "LAN", LOS: time.Date(2028, 12, 28, 0, 45, 16, 0, time.UTC), Outcome: "failed", Reason: "deadline"},
	{PassID: "furrow-a/LAN/21758", Station: "LAN", LOS: time.Date(2028, 12, 28, 1, 57, 49, 0, time.UTC), Outcome: "failed", Reason: "deadline"},
	{PassID: "cirrine-1/LAN/17148", Station: "LAN", LOS: time.Date(2028, 12, 28, 2, 19, 14, 0, time.UTC), Outcome: "failed", Reason: "deadline"},
	{PassID: "furrow-b/LAN/21756", Station: "LAN", LOS: time.Date(2028, 12, 28, 2, 44, 57, 0, time.UTC), Outcome: "failed", Reason: "deadline"},
	{PassID: "brine-3/LAN/15152", Station: "LAN", LOS: time.Date(2028, 12, 28, 7, 26, 4, 0, time.UTC), Outcome: "flown"},
	{PassID: "pipit-1/LAN/14038", Station: "LAN", LOS: time.Date(2028, 12, 28, 8, 12, 34, 0, time.UTC), Outcome: "flown"},
}

// rota is the season's on-call rota, oldest first: four engineers, a week each from Monday
// 08:00 UTC, the weeks of badge No. 039 written as its number.
var rota = []Shift{
	{From: time.Date(2028, 10, 30, 8, 0, 0, 0, time.UTC), Name: "Mirela Dunn"},
	{From: time.Date(2028, 11, 6, 8, 0, 0, 0, time.UTC), Name: "No. 039"},
	{From: time.Date(2028, 11, 13, 8, 0, 0, 0, time.UTC), Name: "Nkem Halloran"},
	{From: time.Date(2028, 11, 20, 8, 0, 0, 0, time.UTC), Name: "Mirela Dunn"},
	{From: time.Date(2028, 11, 27, 8, 0, 0, 0, time.UTC), Name: "Tomás Ferreira-Lund"},
	{From: time.Date(2028, 12, 4, 8, 0, 0, 0, time.UTC), Name: "No. 039"},
	{From: time.Date(2028, 12, 11, 8, 0, 0, 0, time.UTC), Name: "Nkem Halloran"},
	{From: time.Date(2028, 12, 18, 8, 0, 0, 0, time.UTC), Name: "Mirela Dunn"},
	{From: time.Date(2028, 12, 25, 8, 0, 0, 0, time.UTC), Name: "Tomás Ferreira-Lund"},
	{From: time.Date(2029, 1, 1, 8, 0, 0, 0, time.UTC), Name: "No. 039"},
	{From: time.Date(2029, 1, 8, 8, 0, 0, 0, time.UTC), Name: "Nkem Halloran"},
	{From: time.Date(2029, 1, 15, 8, 0, 0, 0, time.UTC), Name: "Mirela Dunn"},
	{From: time.Date(2029, 1, 22, 8, 0, 0, 0, time.UTC), Name: "Tomás Ferreira-Lund"},
	{From: time.Date(2029, 1, 29, 8, 0, 0, 0, time.UTC), Name: "No. 039"},
	{From: time.Date(2029, 2, 5, 8, 0, 0, 0, time.UTC), Name: "Nkem Halloran"},
	{From: time.Date(2029, 2, 12, 8, 0, 0, 0, time.UTC), Name: "Mirela Dunn"},
	{From: time.Date(2029, 2, 19, 8, 0, 0, 0, time.UTC), Name: "Tomás Ferreira-Lund"},
}

// replay runs the night through a pager holding rules, as the new pager will run every
// station: every minute on the minute from 27 December 19:00 to 28 December 09:00 it gives
// the pager each rule's value at Lantern Atoll, whose disk stayed about two thirds full, and
// it returns the alerts the pager sent.
func replay(rules []Rule) []Alert {
	p := &Pager{Rules: rules, Rota: rota}
	start := time.Date(2028, 12, 27, 19, 0, 0, 0, time.UTC)
	var sent []Alert
	for now := start; !now.After(start.Add(14 * time.Hour)); now = now.Add(time.Minute) {
		for _, m := range []struct {
			rule  string
			value float64
		}{
			{"pass-failures", float64(failedWithin(night, "LAN", now, 6*time.Hour))},
			{"failed-pass", float64(failedWithin(night, "LAN", now, time.Minute))},
			{"disk-full", 0.67},
		} {
			if a, ok := p.Observe(m.rule, Sample{At: now, Station: "LAN", Value: m.value}); ok {
				sent = append(sent, a)
			}
		}
	}
	return sent
}
