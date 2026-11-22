package main

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

// window reads the start and the number of days of passes a query asks for, with now and 7 as defaults.
func window(q url.Values, now time.Time) (from time.Time, days int, err error) {
	from, days = now, 7
	if v := q.Get("from"); v != "" {
		if from, err = time.Parse(time.RFC3339, v); err != nil {
			return time.Time{}, 0, errors.New("from must be an RFC 3339 time")
		}
	}
	if v := q.Get("days"); v != "" {
		if days, err = strconv.Atoi(v); err != nil || days < 1 || days > 14 {
			return time.Time{}, 0, errors.New("days must be a whole number from 1 to 14")
		}
	}
	return from, days, nil
}
