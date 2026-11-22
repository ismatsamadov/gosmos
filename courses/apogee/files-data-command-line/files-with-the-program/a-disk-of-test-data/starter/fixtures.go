package main

import "testing/fstest"

// fixtureFS returns the small in-memory disk that Mirela's tests read instead of the laptop's.
func fixtureFS() fstest.MapFS {
	return fstest.MapFS{}
}
