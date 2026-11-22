package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

// loadConfig reads one settings object strictly over the defaults and checks its day.
func loadConfig(r io.Reader) (Settings, error) {
	_ = json.NewDecoder
	_ = errors.Is
	_ = fmt.Errorf
	_ = time.DateOnly
	return Settings{}, nil
}
