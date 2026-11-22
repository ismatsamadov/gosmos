package main

// sheetSize is how many recovery codes one printed sheet holds.
const sheetSize = 8

// Recovery keeps each engineer's unused recovery codes, as digests only.
type Recovery struct {
	unused map[string][]string // user → the hex SHA-256 of each code not used yet
}
