package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/fstest"
	"time"
)

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

// lines is a log that many goroutines may write while the test reads it.
type lines struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// Write adds p to the log.
func (l *lines) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.Write(p)
}

// String returns everything logged so far.
func (l *lines) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.String()
}

// testLinkKey signs download links in tests: 32 bytes, never a production key.
const testLinkKey = "bookingd test link key, 32 bytes"

// plans are the customers' plans: bookings a month. Apogee's own have none.
var plans = map[string]int{"furrowlight": 60, "cirrine": 40, "brinewatch": 12}

// image stands in for Furrow-A's image from Tuesday's 14:02 pass, 10,912,500,000 bytes in life.
func image() []byte {
	data := make([]byte, 64<<10)
	for i := range data {
		data[i] = byte(i*7 + i/256)
	}
	return data
}

// weekCatalog returns the passes the catalog offers on Friday afternoon, by AOS then ID.
func weekCatalog() []PassView {
	pass := func(id, sat, station, aos, los string, maxElevDeg float64) PassView {
		return PassView{ID: id, Pass: Pass{SatID: sat, Station: station, AOS: at(aos), LOS: at(los), MaxElevDeg: maxElevDeg}}
	}
	return []PassView{
		pass("furrow-b/YST/16776", "furrow-b", "YST", "2028-02-04T16:27:34Z", "2028-02-04T16:36:45Z", 68.7),
		pass("cirrine-1/YST/12265", "cirrine-1", "YST", "2028-02-04T17:05:14Z", "2028-02-04T17:15:28Z", 58.9),
		pass("brine-3/PEL/10140", "brine-3", "PEL", "2028-02-04T21:06:32Z", "2028-02-04T21:14:28Z", 20.7),
		pass("brine-3/PEL/10141", "brine-3", "PEL", "2028-02-04T22:43:27Z", "2028-02-04T22:52:45Z", 78.2),
		pass("furrow-a/SQT/16792", "furrow-a", "SQT", "2028-02-05T14:16:50Z", "2028-02-05T14:26:55Z", 84.9),
	}
}

// seed writes Friday afternoon's data folder: the bookings, February's count, Tuesday's landed image and the folder command files go in.
func seed(t *testing.T, dir string) {
	t.Helper()
	bookings := []Booking{
		{ID: "bk-4912", SatID: "furrow-a", Station: "SQT", AOS: at("2028-02-01T14:02:10Z"), LOS: at("2028-02-01T14:11:52Z"), Note: "launch day"},
		{ID: "bk-4913", SatID: "brine-3", Station: "LAN", AOS: at("2028-02-02T12:27:50Z"), LOS: at("2028-02-02T12:34:37Z")},
		{ID: "bk-4933", SatID: "furrow-a", Station: "SQT", AOS: at("2028-02-05T14:16:50Z"), LOS: at("2028-02-05T14:26:55Z")},
		{ID: "bk-4936", SatID: "brine-3", Station: "CWR", AOS: at("2028-02-04T15:43:26Z"), LOS: at("2028-02-04T15:52:39Z")},
	}
	data, err := json.Marshal(bookings)
	if err != nil {
		t.Fatal(err)
	}
	usage := `{"month":"2028-02-01T00:00:00Z","used":{"brinewatch":11,"cirrine":4,"furrowlight":9}}`
	landed := filepath.Join(dir, "deliveries", "bk-4912", "furrow-a-20280201-1402.tif")
	for _, d := range []string{filepath.Join(dir, "commands"), filepath.Dir(landed)} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatal(err)
		}
	}
	for name, content := range map[string][]byte{"bookings.json": data, "usage.json": []byte(usage)} {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(landed, image(), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(landed, at("2028-02-01T14:13:02Z"), at("2028-02-01T14:13:02Z")); err != nil {
		t.Fatal(err)
	}
}

// do sends one request with these headers and returns the answer, its body read and closed.
func do(t *testing.T, method, url string, body io.Reader, header ...string) (*http.Response, []byte) {
	t.Helper()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i+1 < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%s %s: reading the answer: %v", method, url, err)
	}
	return resp, data
}

// bearer gets a token for apiKey from the service at base and returns "Bearer <token>".
func bearer(t *testing.T, base, apiKey string) string {
	t.Helper()
	resp, body := do(t, http.MethodPost, base+"/tokens", nil, "Authorization", "Key "+apiKey)
	var tok TokenResponse
	if err := json.Unmarshal(body, &tok); resp.StatusCode != http.StatusOK || err != nil {
		t.Fatalf("POST /tokens = %d %s, want 200 and a token", resp.StatusCode, body)
	}
	return "Bearer " + tok.Token
}

// form returns a multipart form holding one file part named "file", and its Content-Type.
func form(t *testing.T, filename, partType string, data []byte) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(map[string][]string)
	h["Content-Disposition"] = []string{`form-data; name="file"; filename="` + filename + `"`}
	h["Content-Type"] = []string{partType}
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}
	part.Write(data)
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return &buf, w.FormDataContentType()
}

// errorBody decodes an error answer.
func errorBody(t *testing.T, body []byte) ErrorBody {
	t.Helper()
	var e ErrorBody
	if err := json.Unmarshal(body, &e); err != nil {
		t.Fatalf("error body %s doesn't decode: %v", body, err)
	}
	return e
}

func TestFirstWeek(t *testing.T) {
	dir := t.TempDir()
	seed(t, dir)
	clk := &clock{now: at("2028-02-04T16:00:00Z")}
	var reached atomic.Int32 // requests that reached the weather service
	weather := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached.Add(1)
		http.NotFound(w, r)
	}))
	defer weather.Close()
	customers := testCustomers()
	for i, c := range customers {
		if c.ID == "brinewatch" {
			customers[i].WebhookURL = weather.URL + "/apogee-events" // Tuesday's inside address, replayed to try the guard
		}
	}
	app, err := newApp(Config{
		Dir: dir, Stations: testStations(), Catalog: weekCatalog(),
		Guide:     fstest.MapFS{"index.html": {Data: []byte("<h1>Customer guide</h1>\n")}},
		Customers: customers, TokenKey: []byte(testTokenKey), Now: clk.Now, Rand: rand.New(rand.NewPCG(1, 2)),
	})
	if err != nil {
		t.Fatalf("newApp = %v, want nil", err)
	}
	quota, err := openQuota(filepath.Join(dir, "usage.json"), plans, clk.Now)
	if err != nil {
		t.Fatalf("openQuota = %v, want nil", err)
	}
	files := &Files{Store: app.Service.Store, Commands: filepath.Join(dir, "commands"), Landed: filepath.Join(dir, "deliveries"), LinkKey: []byte(testLinkKey), Now: clk.Now}
	var logs lines
	srv := httptest.NewServer(firstWeek(app, files, quota, log.New(&logs, "", 0)))
	defer srv.Close()
	base := srv.URL
	furrowlight := bearer(t, base, furrowlightKey)

	t.Run("command files", func(t *testing.T) {
		commands := []byte("PING\nDUMP_FILE field-12\n")
		body, ctype := form(t, "commands.txt", "text/plain", commands)
		if resp, got := do(t, http.MethodPost, base+"/bookings/bk-4933/commands", body, "Content-Type", ctype); resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("POST /bookings/bk-4933/commands without a token = %d %s, want 401: a command file comes only from the booking's customer", resp.StatusCode, bytes.TrimSpace(got))
		}
		body, ctype = form(t, "commands.txt", "text/plain", commands)
		resp, got := do(t, http.MethodPost, base+"/bookings/bk-4933/commands", body, "Content-Type", ctype, "Authorization", furrowlight)
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Furrowlight's POST /bookings/bk-4933/commands = %d %s, want 204", resp.StatusCode, bytes.TrimSpace(got))
		}
		if resp.Header.Get("X-Request-Id") == "" {
			t.Errorf("POST /bookings/bk-4933/commands has no X-Request-Id, want the service's middleware around the new routes too")
		}
		if kept, err := os.ReadFile(filepath.Join(dir, "commands", "bk-4933.txt")); err != nil || !bytes.Equal(kept, commands) {
			t.Errorf("commands/bk-4933.txt = %q, %v, want %q", kept, err, commands)
		}
		body, ctype = form(t, "commands.txt", "text/plain", []byte("%PDF-1.7\n"))
		resp, got = do(t, http.MethodPost, base+"/bookings/bk-4933/commands", body, "Content-Type", ctype, "Authorization", furrowlight)
		if e := errorBody(t, got); resp.StatusCode != http.StatusUnsupportedMediaType || e.Message != "command file must be plain text, not application/pdf" {
			t.Errorf("POST /bookings/bk-4933/commands with a PDF = %d %+v, want 415 and the sniffed type", resp.StatusCode, e)
		}
	})

	file := image()
	t.Run("downloads and links", func(t *testing.T) {
		if resp, _ := do(t, http.MethodGet, base+"/bookings/bk-4912/file", nil); resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("GET /bookings/bk-4912/file without a token = %d, want 401", resp.StatusCode)
		}
		resp, got := do(t, http.MethodGet, base+"/bookings/bk-4912/file", nil, "Authorization", furrowlight, "Range", "bytes=60000-")
		if resp.StatusCode != http.StatusPartialContent || !bytes.Equal(got, file[60000:]) {
			t.Errorf("Furrowlight's GET /bookings/bk-4912/file from byte 60000 = %d and %d bytes, want 206 and the last %d", resp.StatusCode, len(got), len(file)-60000)
		}
		if resp, _ := do(t, http.MethodPost, base+"/bookings/bk-4912/link", nil); resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("POST /bookings/bk-4912/link without a token = %d, want 401: only the booking's customer makes links", resp.StatusCode)
		}
		resp, got = do(t, http.MethodPost, base+"/bookings/bk-4912/link", nil, "Authorization", furrowlight)
		var link Link
		if err := json.Unmarshal(got, &link); resp.StatusCode != http.StatusOK || err != nil || !strings.HasPrefix(link.URL, "/files/bk-4912?expires=1833379200&sig=") {
			t.Fatalf("Furrowlight's POST /bookings/bk-4912/link = %d %s, want 200 and a link until Saturday 16:00", resp.StatusCode, bytes.TrimSpace(got))
		}
		resp, got = do(t, http.MethodGet, base+link.URL, nil) // the field team has no token
		if resp.StatusCode != http.StatusOK || !bytes.Equal(got, file) {
			t.Errorf("GET %s with no token = %d and %d bytes, want 200 and the whole image: a link is its own credential", link.URL, resp.StatusCode, len(got))
		}
		resp, got = do(t, http.MethodGet, base+link.URL, nil, "Range", "bytes=65000-")
		if resp.StatusCode != http.StatusPartialContent || !bytes.Equal(got, file[65000:]) {
			t.Errorf("GET of the link from byte 65000 = %d and %d bytes, want 206 and the rest", resp.StatusCode, len(got))
		}
		later := strings.Replace(link.URL, "expires=1833379200", "expires=1833984000", 1)
		if resp, got := do(t, http.MethodGet, base+later, nil); resp.StatusCode != http.StatusForbidden || errorBody(t, got).Message != "invalid link" {
			t.Errorf("GET of the link with its expiry moved a week on = %d %s, want 403 invalid link", resp.StatusCode, bytes.TrimSpace(got))
		}
	})

	t.Run("redirects stay on the site", func(t *testing.T) {
		for next, want := range map[string]string{"https://apogee-bookings.example/renew": "/guide/", "/guide/files.html": "/guide/files.html"} {
			req, err := http.NewRequest(http.MethodGet, base+"/go?next="+next, nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := srv.Client().Transport.RoundTrip(req)
			if err != nil {
				t.Fatalf("GET /go?next=%s: %v", next, err)
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusSeeOther || resp.Header.Get("Location") != want {
				t.Errorf("GET /go?next=%s = %d to %q, want 303 to %q, with no token", next, resp.StatusCode, resp.Header.Get("Location"), want)
			}
		}
	})

	t.Run("plans", func(t *testing.T) {
		brinewatch := bearer(t, base, brinewatchKey)
		twelfth, _ := json.Marshal(Booking{SatID: "brine-3", Station: "PEL", AOS: at("2028-02-04T21:06:32Z"), LOS: at("2028-02-04T21:14:28Z")})
		resp, got := do(t, http.MethodPost, base+"/bookings", bytes.NewReader(twelfth), "Content-Type", "application/json", "Authorization", brinewatch)
		if resp.StatusCode != http.StatusCreated || !bytes.Contains(got, []byte(`"bk-4937"`)) {
			t.Fatalf("Brinewatch's 12th February booking = %d %s, want 201 and bk-4937", resp.StatusCode, bytes.TrimSpace(got))
		}
		thirteenth, _ := json.Marshal(Booking{SatID: "brine-3", Station: "PEL", AOS: at("2028-02-04T22:43:27Z"), LOS: at("2028-02-04T22:52:45Z")})
		resp, got = do(t, http.MethodPost, base+"/bookings", bytes.NewReader(thirteenth), "Content-Type", "application/json", "Authorization", brinewatch)
		e := errorBody(t, got)
		want := "plan of 12 bookings a month used up: the count starts again at 2028-03-01T00:00:00Z"
		if resp.StatusCode != http.StatusForbidden || e.Code != "quota_exceeded" || e.Message != want {
			t.Errorf("Brinewatch's 13th February booking = %d %s, want 403 quota_exceeded %q", resp.StatusCode, bytes.TrimSpace(got), want)
		}
		saved, err := os.ReadFile(filepath.Join(dir, "usage.json"))
		if err != nil || !bytes.Contains(saved, []byte(`"brinewatch":12`)) {
			t.Errorf("usage.json = %s, %v, want Brinewatch at 12", saved, err)
		}
	})

	t.Run("webhooks stay outside", func(t *testing.T) {
		waits := 0
		app.Notifier.Sender.Wait = func(ctx context.Context, d time.Duration) error {
			waits++
			return nil
		}
		d := &Delivery{Customer: "brinewatch", Event: Event{ID: "evt-000131", Type: "data.landed", BookingID: "bk-4936", SatID: "brine-3", Station: "CWR",
			File: "brine-3-20280204-1543.dat", Bytes: 52000000, LandedAt: at("2028-02-04T15:53:41Z")}}
		err := app.Notifier.Sender.Deliver(context.Background(), d)
		if n := reached.Load(); n != 0 {
			t.Fatalf("requests from Brinewatch's event to the weather service at 127.0.0.1 = %d, want 0: it is refused before any connection", n)
		}
		if !errors.Is(err, ErrUndelivered) || !errors.Is(err, errInside) {
			t.Errorf("Deliver(Brinewatch's event to 127.0.0.1) = %v, want ErrUndelivered and errInside", err)
		}
	})

	if resp, _ := do(t, http.MethodGet, base+"/stations", nil); resp.StatusCode != http.StatusOK {
		t.Errorf("GET /stations = %d, want 200: the API is still served beside the new routes", resp.StatusCode)
	}
	if !strings.Contains(logs.String(), " POST /bookings/bk-4933/commands 204 ") {
		t.Errorf("access log = %q, want a line for Furrowlight's upload", logs.String())
	}
}
