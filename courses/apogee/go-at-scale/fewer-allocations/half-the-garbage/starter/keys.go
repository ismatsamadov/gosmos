package main

// passKey is a pass's identity without formatting its ID.
type passKey struct {
	SatID, Station string
	AOS            int64 // Unix seconds
}
