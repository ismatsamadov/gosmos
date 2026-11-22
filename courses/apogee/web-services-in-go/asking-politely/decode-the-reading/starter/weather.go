package main

import (
	"context"
	"fmt"
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

// Reading asks the weather service for the latest reading at the station with this code.
func (c *WeatherClient) Reading(ctx context.Context, code string) (Weather, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/weather/"+code, nil)
	if err != nil {
		return Weather{}, fmt.Errorf("weather for %s: %w", code, err)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return Weather{}, fmt.Errorf("weather for %s: %w", code, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Weather{}, fmt.Errorf("weather for %s: %s", code, resp.Status)
	}
	w, err := decodeWeather(resp.Body)
	if err != nil {
		return Weather{}, fmt.Errorf("weather for %s: %w", code, err)
	}
	return w, nil
}
