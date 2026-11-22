package main

import "testing"

func TestEncodeConfirmation(t *testing.T) {
	for _, tc := range []struct {
		name string
		c    Confirmation
		want string
	}{
		{"confirmed, no note", Confirmation{ID: "bk-2042", Confirmed: true}, `{"id":"bk-2042","confirmed":true}`},
		{"confirmed, with a note", Confirmation{ID: "bk-2041", Confirmed: true, Note: "priority"}, `{"id":"bk-2041","confirmed":true,"note":"priority"}`},
		{"not confirmed yet", Confirmation{ID: "bk-2045"}, `{"id":"bk-2045","confirmed":false}`},
		{"not confirmed, with a note", Confirmation{ID: "bk-2044", Note: "waiting for Paz"}, `{"id":"bk-2044","confirmed":false,"note":"waiting for Paz"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := encodeConfirmation(tc.c)
			if string(got) != tc.want || err != nil {
				t.Errorf("encodeConfirmation(%+v) = %s, %v, want %s, nil", tc.c, got, err, tc.want)
			}
		})
	}
}
