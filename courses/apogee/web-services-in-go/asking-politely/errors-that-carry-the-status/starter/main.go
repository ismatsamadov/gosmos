package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.LimitReader, strings.TrimSpace // TODO: wrap a *StatusError instead of plain text
		return Weather{}, fmt.Errorf("weather for %s: %s", code, resp.Status)
	}
	w, err := decodeWeather(resp.Body)
	if err != nil {
		return Weather{}, fmt.Errorf("weather for %s: %w", code, err)
	}
	return w, nil
}
