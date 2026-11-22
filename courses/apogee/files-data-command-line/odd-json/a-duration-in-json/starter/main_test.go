package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		d    time.Duration
		want string
	}{
		{30 * time.Second, `"30s"`},
		{2 * time.Minute, `"2m0s"`},
		{1500 * time.Millisecond, `"1.5s"`},
	} {
		got, err := json.Marshal(Delay(tc.d))
		if string(got) != tc.want || err != nil {
			t.Errorf("json.Marshal(Delay(%v)) = %s, %v, want %s, nil", tc.d, got, err, tc.want)
		}
	}
	got, err := saveSettings(Settings{Station: "SQT", Day: "2027-08-26", MinElevDeg: 10, Timeout: Delay(30 * time.Second)})
	want := `{"station":"SQT","day":"2027-08-26","min_elev_deg":10,"timeout":"30s"}`
	if string(got) != want || err != nil {
		t.Errorf("saveSettings(Paz's settings) = %s, %v, want %s, nil", got, err, want)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		name, text string
		want       time.Duration
	}{
		{"seconds", `"30s"`, 30 * time.Second},
		{"minutes", `"2m"`, 2 * time.Minute},
		{"both", `"1m30s"`, 90 * time.Second},
		{"escaped digits", "\"\\u0033\\u0030s\"", 30 * time.Second},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var d Delay
			err := json.Unmarshal([]byte(tc.text), &d)
			if time.Duration(d) != tc.want || err != nil {
				t.Errorf("json.Unmarshal(%s) = %v, %v, want %v, nil", tc.text, time.Duration(d), err, tc.want)
			}
		})
	}
	for _, text := range []string{`"thirty"`, `"30"`, `30000000000`, `true`, `["30s"]`} {
		var d Delay
		if err := json.Unmarshal([]byte(text), &d); err == nil {
			t.Errorf("json.Unmarshal(%s) error = nil, want an error", text)
		}
	}
}

func TestDelayRoundTrip(t *testing.T) {
	file := `{"station":"SQT","day":"2027-08-26","min_elev_deg":10,"timeout":"2m"}`
	s, err := loadSettings([]byte(file))
	if time.Duration(s.Timeout) != 2*time.Minute || err != nil {
		t.Fatalf("loadSettings(file).Timeout = %v, %v, want 2m0s, nil", time.Duration(s.Timeout), err)
	}
	saved, err := saveSettings(s)
	want := `{"station":"SQT","day":"2027-08-26","min_elev_deg":10,"timeout":"2m0s"}`
	if string(saved) != want || err != nil {
		t.Fatalf("saveSettings(loaded) = %s, %v, want %s, nil", saved, err, want)
	}
	again, err := loadSettings(saved)
	if again != s || err != nil {
		t.Errorf("loadSettings(saved) = %+v, %v, want %+v, nil", again, err, s)
	}
}
