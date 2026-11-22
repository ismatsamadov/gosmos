package main

import (
	"net/http"
	"time"
)

// newWeatherClient returns a client for the weather service at baseURL whose calls give up after timeout.
func newWeatherClient(baseURL string, timeout time.Duration) *WeatherClient {
	return &WeatherClient{BaseURL: baseURL, HTTP: &http.Client{Timeout: timeout}}
}
