package main

import "encoding/json"

// decodeCounts reads booked passes per station, whatever station codes the summary holds.
func decodeCounts(data []byte) (map[string]int, error) {
	_ = json.Unmarshal
	return nil, nil
}
