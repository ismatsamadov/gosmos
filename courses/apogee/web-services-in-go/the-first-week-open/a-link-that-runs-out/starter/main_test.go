package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// testLinkKey signs download links in tests: 32 bytes, never a production key.
const testLinkKey = "bookingd test link key, 32 bytes"

// sign is the signature a link to path until expires carries: hex HMAC-SHA256 of "<path>?expires=<expires>".
func sign(path, expires, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(path + "?expires=" + expires))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestCheckLink(t *testing.T) {
	const path, expires = "/files/bk-4912", "1833274800" // Friday 4 February, 11:00:00 UTC
	sig := sign(path, expires, testLinkKey)
	thursday := time.Date(2028, 2, 3, 11, 30, 0, 0, time.UTC)
	cases := []struct {
		name               string
		path, expires, sig string
		now                time.Time
		want               error
	}{
		{"Hollis's link, Thursday 11:30", path, expires, sig, thursday, nil},
		{"one second before it runs out", path, expires, sig, time.Date(2028, 2, 4, 10, 59, 59, 0, time.UTC), nil},
		{"the second it runs out", path, expires, sig, time.Date(2028, 2, 4, 11, 0, 0, 0, time.UTC), errLinkExpired},
		{"on Saturday", path, expires, sig, time.Date(2028, 2, 5, 9, 0, 0, 0, time.UTC), errLinkExpired},
		{"its expiry moved a week on", path, "1833879600", sig, thursday, errBadLink},
		{"its expiry moved a day on, read on Saturday", path, "1833361200", sig, time.Date(2028, 2, 5, 12, 0, 0, 0, time.UTC), errBadLink},
		{"its signature on bk-4913's path", "/files/bk-4913", expires, sig, thursday, errBadLink},
		{"signed with another key", path, expires, sign(path, expires, "someone else's key, 32 bytes!!!"), thursday, errBadLink},
		{"no signature", path, expires, "", thursday, errBadLink},
		{"a signature that isn't hex", path, expires, strings.Repeat("zz", 32), thursday, errBadLink},
		{"an expiry that isn't a number", path, "friday", sig, thursday, errBadLink},
		{"no expiry", path, "", sig, thursday, errBadLink},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := checkLink(c.path, c.expires, c.sig, []byte(testLinkKey), c.now)
			if !errors.Is(err, c.want) || (c.want == nil && err != nil) {
				t.Errorf("checkLink(%q, %q, %.8s…, at %s) = %v, want %v", c.path, c.expires, c.sig, c.now.Format(time.RFC3339), err, c.want)
			}
		})
	}
}

func TestCheckLinkInConstantTime(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "link.go", nil, 0)
	if err != nil {
		t.Fatalf("parser.ParseFile(link.go) error = %v, want nil", err)
	}
	calls := 0
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "checkLink" || fn.Body == nil {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if pkg, ok := sel.X.(*ast.Ident); ok && (pkg.Name == "hmac" && sel.Sel.Name == "Equal" || pkg.Name == "subtle" && sel.Sel.Name == "ConstantTimeCompare") {
						calls++
					}
				}
			}
			return true
		})
	}
	if calls == 0 {
		t.Errorf("checkLink calls hmac.Equal %d times, want the signatures compared in constant time: == stops at the first byte that differs", calls)
	}
}

// clock is a fake clock that moves only when the test moves it.
type clock struct {
	mu  sync.Mutex
	now time.Time
}

// Now returns the clock's time.
func (c *clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Set moves the clock to t.
func (c *clock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t
}

// bookings is a store holding Furrowlight's launch-day booking, enough for the link handlers.
type bookings map[string]Booking

// Create is not used here.
func (s bookings) Create(ctx context.Context, b Booking) (Booking, error) {
	return Booking{}, errors.New("not in this test")
}

// Get returns the booking with this ID, or ErrNotFound.
func (s bookings) Get(ctx context.Context, id string) (Booking, error) {
	b, ok := s[id]
	if !ok {
		return Booking{}, ErrNotFound
	}
	return b, nil
}

// List is not used here.
func (s bookings) List(ctx context.Context) ([]Booking, error) {
	return nil, errors.New("not in this test")
}

// Delete is not used here.
func (s bookings) Delete(ctx context.Context, id string) error { return errors.New("not in this test") }

// get sends a GET to url with these headers, and returns the status and the body.
func get(t *testing.T, srv *httptest.Server, url string, header ...string) (int, []byte) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, srv.URL+url, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i+1 < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("GET %s: reading the body: %v", url, err)
	}
	return resp.StatusCode, body
}

func TestTheFieldTeamsLink(t *testing.T) {
	landed := t.TempDir()
	image := bytes.Repeat([]byte("Furrow-A over Salar Quieto, 1 February 14:02 "), 1000)
	if err := os.Mkdir(filepath.Join(landed, "bk-4912"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(landed, "bk-4912", "furrow-a-20280201-1402.tif"), image, 0600); err != nil {
		t.Fatal(err)
	}
	clk := &clock{now: time.Date(2028, 2, 3, 11, 0, 0, 0, time.UTC)}
	f := &Files{
		Store:   bookings{"bk-4912": {ID: "bk-4912", SatID: "furrow-a", Station: "SQT", AOS: time.Date(2028, 2, 1, 14, 2, 10, 0, time.UTC), LOS: time.Date(2028, 2, 1, 14, 11, 52, 0, time.UTC)}},
		Landed:  landed,
		LinkKey: []byte(testLinkKey),
		Now:     clk.Now,
	}
	furrowlight := Customer{ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /bookings/{id}/link", func(w http.ResponseWriter, r *http.Request) {
		f.link(w, r.WithContext(withCustomer(r.Context(), furrowlight)))
	})
	mux.HandleFunc("GET /files/{id}", f.shared) // the field team has no token
	srv := httptest.NewServer(mux)
	defer srv.Close()

	resp, err := srv.Client().Post(srv.URL+"/bookings/bk-4912/link", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	var link Link
	err = json.NewDecoder(resp.Body).Decode(&link)
	resp.Body.Close()
	if want := "/files/bk-4912?expires=1833274800&sig=" + sign("/files/bk-4912", "1833274800", testLinkKey); err != nil || link.URL != want {
		t.Fatalf("Hollis's POST /bookings/bk-4912/link = %+v, %v, want the link %s", link, err, want)
	}
	if status, body := get(t, srv, link.URL); status != http.StatusOK || !bytes.Equal(body, image) {
		t.Errorf("GET of Hollis's link on Thursday = %d and %d bytes, want 200 and the whole image", status, len(body))
	}
	if status, body := get(t, srv, link.URL, "Range", "bytes=40000-"); status != http.StatusPartialContent || !bytes.Equal(body, image[40000:]) {
		t.Errorf("GET of Hollis's link from byte 40000 = %d and %d bytes, want 206 and the rest", status, len(body))
	}
	longer := strings.Replace(link.URL, "expires=1833274800", "expires=1833879600", 1)
	if status, body := get(t, srv, longer); status != http.StatusForbidden || !bytes.Contains(body, []byte("invalid link")) {
		t.Errorf("GET of the link with its expiry moved a week on = %d %.60q, want 403 invalid link", status, body)
	}
	clk.Set(time.Date(2028, 2, 4, 11, 0, 0, 0, time.UTC))
	if status, body := get(t, srv, link.URL); status != http.StatusForbidden || !bytes.Contains(body, []byte("link expired")) {
		t.Errorf("GET of Hollis's link on Friday at 11:00:00 = %d %.60q, want 403 link expired", status, body)
	}
}
