package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// windCall asks the weather service about one station and returns the poll's wind call for it.
func windCall(ctx context.Context, c *WeatherClient, code string, now time.Time) string {
	_, _ = errors.AsType[*StatusError], fmt.Sprintf
	return ""
}
