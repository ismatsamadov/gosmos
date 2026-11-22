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
		start, _ := io.ReadAll(io.LimitReader(resp.Body, 512)) // the status matters even if the body breaks off
		return Weather{}, fmt.Errorf("weather for %s: %w", code, &StatusError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       strings.TrimSpace(string(start)),
		})
	}
	w, err := decodeWeather(resp.Body)
	if err != nil {
		return Weather{}, fmt.Errorf("weather for %s: %w", code, err)
	}
	return w, nil
}
