package main

import (
	"strings"
	"testing"
	"time"
)

func TestDecodeWeather(t *testing.T) {
	ystmark := Weather{Station: "YST", At: time.Date(2027, 11, 11, 3, 10, 0, 0, time.UTC), WindKmh: 64, GustKmh: 81, TempC: -12.5}
	tests := []struct {
		name    string
		body    string
		want    Weather
		wantErr bool
	}{
		{"Ystmark's whole reading",
			`{"station":"YST","at":"2027-11-11T03:10:00Z","wind_kmh":64,"gust_kmh":81,"temp_c":-12.5}`,
			ystmark, false},
		{"humidity, added last spring",
			`{"station":"YST","at":"2027-11-11T03:10:00Z","wind_kmh":64,"gust_kmh":81,"temp_c":-12.5,"humidity_pct":91}`,
			ystmark, false},
		{"a field nobody has seen yet",
			`{"visibility_km":0.4,"station":"YST","at":"2027-11-11T03:10:00Z","wind_kmh":64,"gust_kmh":81,"temp_c":-12.5}`,
			ystmark, false},
		{"no station", `{"at":"2027-11-11T03:10:00Z","wind_kmh":64}`, Weather{}, true},
		{"no at", `{"station":"LAN","wind_kmh":12}`, Weather{}, true},
		{"broken JSON", `{"station":"LAN","at":`, Weather{}, true},
		{"an empty body", ``, Weather{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeWeather(strings.NewReader(tt.body))
			if tt.wantErr {
				if err == nil {
					t.Errorf("decodeWeather(`%s`) error = nil, want an error", tt.body)
				}
				return
			}
			if err != nil {
				t.Fatalf("decodeWeather(`%s`) error = %v, want nil", tt.body, err)
			}
			if got.Station != tt.want.Station || !got.At.Equal(tt.want.At) || got.WindKmh != tt.want.WindKmh ||
				got.GustKmh != tt.want.GustKmh || got.TempC != tt.want.TempC {
				t.Errorf("decodeWeather(`%s`) = %+v, want %+v", tt.body, got, tt.want)
			}
		})
	}
}
