package main

// loadRegistry returns the old fleet list: each satellite's operator, by satellite ID.
func loadRegistry() map[string]string {
	return map[string]string{
		"pipit-1":   "Apogee",
		"pipit-2":   "Apogee",
		"pipit-3":   "Apogee",
		"furrow-a":  "Furrowlight",
		"furrow-b":  "Furrowlight",
		"brine-2":   "Brinewatch",
		"brine-3":   "Brinewatch",
		"thistle":   "Kessick Polytechnic",
		"cirrine-1": "Cirrine Atmospherics",
	}
}
