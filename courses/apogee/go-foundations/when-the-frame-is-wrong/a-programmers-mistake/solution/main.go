package main

// requireLimit rejects an impossible internal decoder configuration.
func requireLimit(limit int) {
	if limit <= 0 {
		panic("invalid decoder limit")
	}
}
