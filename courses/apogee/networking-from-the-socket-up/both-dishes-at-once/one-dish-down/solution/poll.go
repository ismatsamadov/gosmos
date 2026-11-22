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
	results := make([]Result, len(clients))
	var wg sync.WaitGroup
	for i, c := range clients {
		wg.Go(func() {
			st, err := c.status()
			results[i] = Result{Name: c.Name, Status: st, Err: err}
		})
	}
	wg.Wait()
	return results
}
