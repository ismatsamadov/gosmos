package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// landedAt is when Furrow-A's image from Tuesday's 14:02 pass landed.
var landedAt = time.Date(2028, 2, 1, 14, 13, 2, 0, time.UTC)

// image stands in for that image, 10,912,500,000 bytes in life.
func image() []byte {
	data := make([]byte, 64<<10)
	for i := range data {
		data[i] = byte(i*7 + i/256)
	}
	return data
}

// bookings is a store holding Furrowlight's bookings, enough for the download handler.
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

// serve starts a test server whose GET /bookings/{id}/file downloads as Furrowlight, from a delivery folder holding Tuesday's image.
func serve(t *testing.T) *httptest.Server {
	t.Helper()
	landed := t.TempDir()
	path := filepath.Join(landed, "bk-4912", "furrow-a-20280201-1402.tif")
	if err := os.Mkdir(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, image(), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, landedAt, landedAt); err != nil {
		t.Fatal(err)
	}
	f := &Files{
		Store: bookings{
			"bk-4912": {ID: "bk-4912", SatID: "furrow-a", Station: "SQT", AOS: time.Date(2028, 2, 1, 14, 2, 10, 0, time.UTC), LOS: time.Date(2028, 2, 1, 14, 11, 52, 0, time.UTC)},
			"bk-4921": {ID: "bk-4921", SatID: "furrow-a", Station: "SQT", AOS: time.Date(2028, 2, 4, 14, 36, 38, 0, time.UTC), LOS: time.Date(2028, 2, 4, 14, 46, 32, 0, time.UTC)},
		},
		Landed: landed,
	}
	furrowlight := Customer{ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /bookings/{id}/file", func(w http.ResponseWriter, r *http.Request) {
		f.download(w, r.WithContext(withCustomer(r.Context(), furrowlight)))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// get asks srv for a file with these headers, and returns the answer and its body.
func get(t *testing.T, srv *httptest.Server, path string, header ...string) (*http.Response, []byte) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, srv.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i+1 < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("GET %s: reading the body: %v", path, err)
	}
	return resp, body
}

func TestWholeImage(t *testing.T) {
	srv := serve(t)
	resp, body := get(t, srv, "/bookings/bk-4912/file")
	if resp.StatusCode != http.StatusOK || !bytes.Equal(body, image()) {
		t.Fatalf("GET /bookings/bk-4912/file = %d and %d bytes, want 200 and all 65536", resp.StatusCode, len(body))
	}
	if got := resp.Header.Get("Accept-Ranges"); got != "bytes" {
		t.Errorf("Accept-Ranges = %q, want %q: it tells clients they can resume", got, "bytes")
	}
	if got, want := resp.Header.Get("Last-Modified"), "Tue, 01 Feb 2028 14:13:02 GMT"; got != want {
		t.Errorf("Last-Modified = %q, want %q, the time the image landed", got, want)
	}
}

func TestResume(t *testing.T) {
	srv := serve(t)
	have := image()[:40000] // where Hollis's fourth try broke
	resp, rest := get(t, srv, "/bookings/bk-4912/file", "Range", "bytes=40000-")
	if resp.StatusCode != http.StatusPartialContent {
		t.Fatalf("GET with Range: bytes=40000- = %d, want 206 Partial Content", resp.StatusCode)
	}
	if got, want := resp.Header.Get("Content-Range"), "bytes 40000-65535/65536"; got != want {
		t.Errorf("Content-Range = %q, want %q", got, want)
	}
	if whole := append(have, rest...); !bytes.Equal(whole, image()) {
		t.Errorf("the first 40000 bytes and the %d resumed = %d bytes unlike the image, want the image whole", len(rest), len(whole))
	}

	if resp, _ := get(t, srv, "/bookings/bk-4912/file", "Range", "bytes=65536-"); resp.StatusCode != http.StatusRequestedRangeNotSatisfiable {
		t.Errorf("GET with Range: bytes=65536-, a download already whole = %d, want 416", resp.StatusCode)
	}
}

func TestResumeOnlyTheSameImage(t *testing.T) {
	srv := serve(t)
	resp, rest := get(t, srv, "/bookings/bk-4912/file", "Range", "bytes=40000-", "If-Range", "Tue, 01 Feb 2028 14:13:02 GMT")
	if resp.StatusCode != http.StatusPartialContent || len(rest) != 65536-40000 {
		t.Errorf("a resume with If-Range: the image's own Last-Modified = %d and %d bytes, want 206 and the last 25536", resp.StatusCode, len(rest))
	}
	resp, whole := get(t, srv, "/bookings/bk-4912/file", "Range", "bytes=40000-", "If-Range", "Tue, 01 Feb 2028 14:12:58 GMT")
	if resp.StatusCode != http.StatusOK || !bytes.Equal(whole, image()) {
		t.Errorf("a resume with If-Range: an older Last-Modified = %d and %d bytes, want 200 and the whole image: the file has changed since", resp.StatusCode, len(whole))
	}
}

func TestNotLandedYet(t *testing.T) {
	srv := serve(t)
	if resp, body := get(t, srv, "/bookings/bk-4921/file"); resp.StatusCode != http.StatusNotFound || !bytes.Contains(body, []byte("no file for bk-4921 yet")) {
		t.Errorf("GET /bookings/bk-4921/file before Friday's pass = %d %s, want 404 no file for bk-4921 yet", resp.StatusCode, body)
	}
}
