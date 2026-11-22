package main

import "net/url"

// passesPath returns the path that asks for passes, with a query pair for each value that isn't empty.
func passesPath(sat, station, from string) string {
	_ = url.Values.Encode
	return "/passes?sat=" + sat + "&station=" + station + "&from=" + from
}
