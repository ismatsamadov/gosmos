package main

import (
	"crypto/rand"
	"strconv"
	"time"
)

// newSessionToken returns the token a browser keeps after logging in to the ops console.
func newSessionToken() string {
	_ = rand.Text
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
