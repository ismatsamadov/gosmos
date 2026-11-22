package main

import (
	"encoding/json"
	"errors"
	"testing"
)

const bk2041 = `{"id":"bk-2041","sat_id":"furrow-a","station":"SQT","aos":"2027-08-11T14:02:10Z","los":"2027-08-11T14:11:52Z","note":"priority"}`

func TestConfirm(t *testing.T) {
	for _, tc := range []struct {
		name, data, want string
	}{
		{"as Yusra sent it", bk2041, bk2041},
		{"spread over lines", "{\n  \"id\": \"bk-2042\",\n  \"sat_id\": \"furrow-b\",\n  \"station\": \"SQT\",\n  \"aos\": \"2027-08-11T14:49:40Z\",\n  \"los\": \"2027-08-11T14:58:31Z\"\n}\n",
			`{"id":"bk-2042","sat_id":"furrow-b","station":"SQT","aos":"2027-08-11T14:49:40Z","los":"2027-08-11T14:58:31Z"}`},
		{"an empty note", `{"id":"bk-2043","sat_id":"cirrine-1","station":"SQT","aos":"2027-08-11T16:20:05Z","los":"2027-08-11T16:30:41Z","note":""}`,
			`{"id":"bk-2043","sat_id":"cirrine-1","station":"SQT","aos":"2027-08-11T16:20:05Z","los":"2027-08-11T16:30:41Z"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := confirm([]byte(tc.data))
			if string(got) != tc.want || !json.Valid(got) || err != nil {
				t.Errorf("confirm(%s) = %s, %v, want %s, nil", tc.name, got, err, tc.want)
			}
		})
	}
}

func TestConfirmRefuses(t *testing.T) {
	for _, tc := range []struct{ name, data string }{
		{"no id", `{"sat_id":"furrow-a","station":"SQT","aos":"2027-08-11T14:02:10Z","los":"2027-08-11T14:11:52Z"}`},
		{"an empty id", `{"id":"","sat_id":"furrow-a","station":"SQT","aos":"2027-08-11T14:02:10Z","los":"2027-08-11T14:11:52Z"}`},
		{"a null id", `{"id":null,"sat_id":"furrow-a","station":"SQT","aos":"2027-08-11T14:02:10Z","los":"2027-08-11T14:11:52Z"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := confirm([]byte(tc.data))
			if got != nil || err == nil {
				t.Errorf("confirm(%s) = %s, %v, want nil, an error", tc.name, got, err)
			}
		})
	}
	got, err := confirm([]byte(bk2041[:40]))
	var syntax *json.SyntaxError
	if got != nil || !errors.As(err, &syntax) {
		t.Errorf("confirm(cut short) = %s, %v, want nil, a syntax error", got, err)
	}
}
