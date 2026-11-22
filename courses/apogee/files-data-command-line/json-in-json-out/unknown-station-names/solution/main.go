package main

import "encoding/json"

// decodeCounts reads booked passes per station, whatever station codes the summary holds.
func decodeCounts(data []byte) (map[string]int, error) {
	var counts map[string]int
	if err := json.Unmarshal(data, &counts); err != nil {
		return nil, err
	}
	return counts, nil
}
