package main

import (
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

// newWeatherClient returns a client for the weather service at baseURL whose calls give up after timeout.
func newWeatherClient(baseURL string, timeout time.Duration) *WeatherClient {
	return &WeatherClient{BaseURL: baseURL, HTTP: &http.Client{Timeout: timeout}}
}
