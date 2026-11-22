package main

import "encoding/json"

// MaskUpdate is a change to one station's settings, as gctl receives it.
type MaskUpdate struct {
	Station string   `json:"station"`
	MaskDeg *float64 `json:"mask_deg"`
}

// decodeMaskUpdate returns the mask a station has after an update: oldDeg unless the update sets one.
func decodeMaskUpdate(data []byte, oldDeg float64) (float64, error) {
	var u MaskUpdate
	if err := json.Unmarshal(data, &u); err != nil {
		return oldDeg, err
	}
	if u.MaskDeg == nil {
		return oldDeg, nil
	}
	return *u.MaskDeg, nil
}
