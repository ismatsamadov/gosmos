package main

import "time"

// networkTime rewrites an RFC 3339 timestamp the way the network logs it: the same instant, in UTC.
func networkTime(text string) (string, error) {
	t, err := time.Parse(time.RFC3339, text)
	if err != nil {
		return "", err
	}
	return t.UTC().Format(time.RFC3339), nil
}
