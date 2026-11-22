package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// linkLife is how long a download link works after it is made.
const linkLife = 24 * time.Hour

// Link is the body of POST /bookings/{id}/link: a download link that works without a token until it expires.
type Link struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// A download link that fails its check gets one of these answers, and nothing more.
var (
	errBadLink     = &APIError{Status: http.StatusForbidden, Code: "forbidden", Message: "invalid link"}
	errLinkExpired = &APIError{Status: http.StatusForbidden, Code: "forbidden", Message: "link expired"}
)

// linkSig returns the hex HMAC-SHA256, under key, of a link's path and expiry: the link as it reads up to "&sig=".
func linkSig(path string, expires int64, key []byte) string {
	mac := hmac.New(sha256.New, key)
	fmt.Fprintf(mac, "%s?expires=%d", path, expires)
	return hex.EncodeToString(mac.Sum(nil))
}

// linkFor returns a link to path that works until expires, signed with key.
func linkFor(path string, expires time.Time, key []byte) string {
	return fmt.Sprintf("%s?expires=%d&sig=%s", path, expires.Unix(), linkSig(path, expires.Unix(), key))
}

// checkLink returns nil when sig signs path and expires under key and now is before the expiry: errBadLink for a link anyone changed, errLinkExpired for a genuine one past its time.
func checkLink(path, expires, sig string, key []byte, now time.Time) error {
	exp, err := strconv.ParseInt(expires, 10, 64)
	if err != nil || !hmac.Equal([]byte(sig), []byte(linkSig(path, exp, key))) {
		return errBadLink
	}
	if !now.Before(time.Unix(exp, 0)) {
		return errLinkExpired
	}
	return nil
}
