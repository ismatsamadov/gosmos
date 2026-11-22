package main

import "sync"

// networkOnce returns a function that gives every caller the result of one call of load: its first.
func networkOnce(load func() (*Network, error)) func() (*Network, error) {
	return sync.OnceValues(load)
}
