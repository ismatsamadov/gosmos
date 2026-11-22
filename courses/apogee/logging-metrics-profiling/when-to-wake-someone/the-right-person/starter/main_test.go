package main

import (
	"fmt"
	"slices"
	"testing"
	"time"
)

// at returns a moment of the season, in UTC.
func at(year int, month time.Month, day, h, m int) time.Time {
	return time.Date(year, month, day, h, m, 0, 0, time.UTC)
}

// rota is the season's on-call rota, oldest first: four engineers, a week each from Monday
// 08:00 UTC, the weeks of badge No. 039 written as its number.
var rota = []Shift{
	{From: at(2028, 10, 30, 8, 0), Name: "Mirela Dunn"},
	{From: at(2028, 11, 6, 8, 0), Name: "No. 039"},
	{From: at(2028, 11, 13, 8, 0), Name: "Nkem Halloran"},
	{From: at(2028, 11, 20, 8, 0), Name: "Mirela Dunn"},
	{From: at(2028, 11, 27, 8, 0), Name: "Tomás Ferreira-Lund"},
	{From: at(2028, 12, 4, 8, 0), Name: "No. 039"},
	{From: at(2028, 12, 11, 8, 0), Name: "Nkem Halloran"},
	{From: at(2028, 12, 18, 8, 0), Name: "Mirela Dunn"},
	{From: at(2028, 12, 25, 8, 0), Name: "Tomás Ferreira-Lund"},
	{From: at(2029, 1, 1, 8, 0), Name: "No. 039"},
	{From: at(2029, 1, 8, 8, 0), Name: "Nkem Halloran"},
	{From: at(2029, 1, 15, 8, 0), Name: "Mirela Dunn"},
	{From: at(2029, 1, 22, 8, 0), Name: "Tomás Ferreira-Lund"},
	{From: at(2029, 1, 29, 8, 0), Name: "No. 039"},
	{From: at(2029, 2, 5, 8, 0), Name: "Nkem Halloran"},
	{From: at(2029, 2, 12, 8, 0), Name: "Mirela Dunn"},
	{From: at(2029, 2, 19, 8, 0), Name: "Tomás Ferreira-Lund"},
}

// severities name each Severity for messages.
var severities = map[Severity]string{Info: "Info", Warning: "Warning", Critical: "Critical"}

func TestRoute(t *testing.T) {
	tests := []struct {
		name  string
		alert Alert
		sev   Severity
		want  string
	}{
		{"Thursday 28 December 00:45, Tomás's week",
			Alert{At: at(2028, 12, 28, 0, 45), Rule: "pass-failures", Station: "LAN", Kind: "FIRING", Value: 3}, Critical, "Tomás Ferreira-Lund"},
		{"Thursday 28 December 08:45, the RESOLVED goes to Tomás too",
			Alert{At: at(2028, 12, 28, 8, 45), Rule: "pass-failures", Station: "LAN", Kind: "RESOLVED", Value: 0}, Critical, "Tomás Ferreira-Lund"},
		{"Monday 1 January 07:59, still Tomás",
			Alert{At: at(2029, 1, 1, 7, 59), Rule: "pass-failures", Station: "LAN", Kind: "FIRING", Value: 2}, Critical, "Tomás Ferreira-Lund"},
		{"Monday 1 January 08:00, your week",
			Alert{At: at(2029, 1, 1, 8, 0), Rule: "pass-failures", Station: "LAN", Kind: "FIRING", Value: 2}, Critical, "No. 039"},
		{"Thursday 11 January 11:05, Nkem's week",
			Alert{At: at(2029, 1, 11, 11, 5), Rule: "disk-full", Station: "LAN", Kind: "FIRING", Value: 0.92}, Critical, "Nkem Halloran"},
		{"a single failed pass, for the morning list",
			Alert{At: at(2028, 12, 27, 19, 42), Rule: "failed-pass", Station: "LAN", Kind: "FIRING", Value: 1}, Warning, "morning list"},
		{"an Info alert, only logged",
			Alert{At: at(2028, 12, 28, 7, 27), Rule: "pass-flown", Station: "LAN", Kind: "FIRING", Value: 1}, Info, "log"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := slices.Clone(rota)
			got := route(tt.alert, tt.sev, rota)
			call := fmt.Sprintf("route(%s %s, %s, rota)", tt.alert.At.Format("Jan 2"), tt.alert, severities[tt.sev])
			if got.To != tt.want {
				t.Errorf("%s.To = %q, want %q", call, got.To, tt.want)
			}
			if got.To = tt.alert.To; got != tt.alert {
				t.Errorf("%s changed the alert to %s, want only To set", call, got)
			}
			if !slices.Equal(rota, before) {
				t.Errorf("%s changed the rota, want it as it was", call)
			}
		})
	}
}
