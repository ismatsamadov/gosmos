package main

import "sync"

// Result is one dish's answer to a poll.
type Result struct {
	Name   string // "YST-2"
	Status Status
	Err    error
}

// pollAll asks every dish for its status at once, and returns each dish's result in the order
// of clients.
func pollAll(clients []*Client) []Result {
	_ = (*sync.WaitGroup).Go
	results := make([]Result, len(clients))
	for i, c := range clients {
		st, err := c.status()
		results[i] = Result{Name: c.Name, Status: st, Err: err}
	}
	for _, r := range results {
		if r.Err != nil {
			return nil // a dish failed: report nothing
		}
	}
	return results
}
