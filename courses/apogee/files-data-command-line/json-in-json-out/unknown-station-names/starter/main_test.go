package main

import (
	"maps"
	"testing"
)

func TestDecodeCounts(t *testing.T) {
	for _, tc := range []struct {
		name, data string
		want       map[string]int
	}{
		{"last week", `{"YST":6,"SQT":4}`, map[string]int{"YST": 6, "SQT": 4}},
		{"Cape Wren joins", `{"YST":6,"SQT":4,"CWR":2}`, map[string]int{"YST": 6, "SQT": 4, "CWR": 2}},
		{"none at Cape Wren", `{"SQT":4,"CWR":0}`, map[string]int{"SQT": 4, "CWR": 0}},
		{"an empty week", `{}`, map[string]int{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeCounts([]byte(tc.data))
			if !maps.Equal(got, tc.want) || err != nil {
				t.Errorf("decodeCounts(%s) = %v, %v, want %v, nil", tc.data, got, err, tc.want)
			}
		})
	}
}

func TestDecodeCountsInvalid(t *testing.T) {
	for _, tc := range []struct{ name, data string }{
		{"half a pass", `{"SQT":4,"CWR":2.5}`},
		{"a word for a count", `{"SQT":"four"}`},
		{"a list, not an object", `["SQT","CWR"]`},
		{"cut short", `{"SQT":4,`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeCounts([]byte(tc.data))
			if got != nil || err == nil {
				t.Errorf("decodeCounts(%s) = %v, %v, want nil, an error", tc.data, got, err)
			}
		})
	}
}
