package main

// knownSat reports whether id is one of the seven satellites the service books passes for. Thistle re-entered in July, so it is not one of them.
func knownSat(id string) bool {
	switch id {
	case "pipit-1", "pipit-2", "pipit-3", "furrow-a", "furrow-b", "brine-3", "cirrine-1":
		return true
	}
	return false
}

// knownStation reports whether code is exactly one of the five station codes.
func knownStation(code string) bool {
	switch code {
	case "YST", "PEL", "SQT", "LAN", "CWR":
		return true
	}
	return false
}
