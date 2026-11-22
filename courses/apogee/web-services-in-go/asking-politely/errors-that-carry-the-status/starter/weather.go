package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// Weather is one station's latest reading from the weather service.
type Weather struct {
	Station string    `json:"station"`
	At      time.Time `json:"at"`
	WindKmh float64   `json:"wind_kmh"`
	GustKmh float64   `json:"gust_kmh"`
	TempC   float64   `json:"temp_c"`
}

// WeatherClient asks the weather service for station readings.
type WeatherClient struct {
	BaseURL string       // the service's address, in tests an httptest server's URL
	HTTP    *http.Client // with a Timeout
}

// StatusError is an answer from the weather service other than 200 OK.
type StatusError struct {
	StatusCode int    // 503
	Status     string // "503 Service Unavailable"
	Body       string // at most the first 512 bytes it sent, trimmed
}

// Error returns the status, then what the service said.
func (e *StatusError) Error() string { return e.Status + ": " + e.Body }

// newWeatherClient returns a client for the weather service at baseURL whose calls give up after timeout.
func newWeatherClient(baseURL string, timeout time.Duration) *WeatherClient {
	return &WeatherClient{BaseURL: baseURL, HTTP: &http.Client{Timeout: timeout}}
}

// decodeWeather reads one reading from r and refuses one without a station or a time.
func decodeWeather(r io.Reader) (Weather, error) {
	var w Weather
	if err := json.NewDecoder(r).Decode(&w); err != nil {
		return Weather{}, err
	}
	if w.Station == "" {
		return Weather{}, errors.New(`reading has no "station"`)
	}
	if w.At.IsZero() {
		return Weather{}, errors.New(`reading has no "at"`)
	}
	return w, nil
}
