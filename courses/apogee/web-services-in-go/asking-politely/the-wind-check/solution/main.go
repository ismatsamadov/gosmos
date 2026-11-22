package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// windCall asks the weather service about one station and returns the poll's wind call for it.
func windCall(ctx context.Context, c *WeatherClient, code string, now time.Time) string {
	w, err := c.Reading(ctx, code)
	if se, ok := errors.AsType[*StatusError](err); ok && se.StatusCode == 503 {
		return "NO-GO: mast silent"
	}
	if err != nil {
		return "NO-GO: no weather"
	}
	if age := now.Sub(w.At); age > 15*time.Minute {
		return fmt.Sprintf("NO-GO: reading %v old", age)
	}
	switch {
	case w.WindKmh >= 72:
		return fmt.Sprintf("NO-GO: wind %.0f km/h, stow", w.WindKmh)
	case w.WindKmh >= 50:
		return fmt.Sprintf("GO: wind %.0f km/h, warning", w.WindKmh)
	}
	return fmt.Sprintf("GO: wind %.0f km/h", w.WindKmh)
}
