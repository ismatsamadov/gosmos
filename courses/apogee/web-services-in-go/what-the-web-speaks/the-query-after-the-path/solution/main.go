package main

import "net/url"

// passesPath returns the path that asks for passes, with a query pair for each value that isn't empty.
func passesPath(sat, station, from string) string {
	q := url.Values{}
	if sat != "" {
		q.Set("sat", sat)
	}
	if station != "" {
		q.Set("station", station)
	}
	if from != "" {
		q.Set("from", from)
	}
	if len(q) == 0 {
		return "/passes"
	}
	return "/passes?" + q.Encode()
}
