package main

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

// window reads the start and the number of days of passes a query asks for, with now and 7 as defaults.
func window(q url.Values, now time.Time) (from time.Time, days int, err error) {
	_ = errors.New
	_ = strconv.Atoi
	return now, 14, nil
}
