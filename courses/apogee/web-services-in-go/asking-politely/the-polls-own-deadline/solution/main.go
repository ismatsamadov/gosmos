package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

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
	var w Weather
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return Weather{}, fmt.Errorf("weather for %s: %w", code, err)
	}
	resp.Body.Close()
	return w, nil
}
