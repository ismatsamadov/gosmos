package main

import (
	"net/http"
	"time"
)

// cookieName is the name of the cookie that carries a browser's session ID.
const cookieName = "session"

// sessionLife is how long a session lasts after it was last used.
const sessionLife = 30 * time.Minute

// sessionCookie returns the cookie that hands session id to the browser.
func sessionCookie(id string) *http.Cookie {
	return &http.Cookie{
		Name:  cookieName,
		Value: id,
	}
}
