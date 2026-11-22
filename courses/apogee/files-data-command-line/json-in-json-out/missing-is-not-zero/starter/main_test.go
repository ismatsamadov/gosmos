package main

import (
	"math"
	"testing"
)

func TestDecodeMaskUpdate(t *testing.T) {
	for _, tc := range []struct {
		name, data string
		want       float64
	}{
		{"no mask in the update", `{"station":"SQT"}`, 3},
		{"a null mask", `{"station":"SQT","mask_deg":null}`, 3},
		{"down to the horizon", `{"station":"SQT","mask_deg":0}`, 0},
		{"a higher mask", `{"station":"SQT","mask_deg":5}`, 5},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeMaskUpdate([]byte(tc.data), 3)
			if math.Abs(got-tc.want) > 0.01 || err != nil {
				t.Errorf("decodeMaskUpdate(%s, 3) = %.2f, %v, want %.2f, nil", tc.data, got, err, tc.want)
			}
		})
	}
	data := `{"station":"SQT","mask_deg":"low"}`
	got, err := decodeMaskUpdate([]byte(data), 3)
	if got != 3 || err == nil {
		t.Errorf("decodeMaskUpdate(%s, 3) = %.2f, %v, want 3.00, an error", data, got, err)
	}
}
