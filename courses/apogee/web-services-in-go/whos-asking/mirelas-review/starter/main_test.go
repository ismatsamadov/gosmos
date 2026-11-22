package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"
)

// tokenKey is the test token key: 32 bytes, like the service's real key, which never appears in a lesson.
const tokenKey = "bookingd test token key 32 bytes"

// clock is the service's clock in these tests: it moves only when the review calls advance.
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

// Add moves the clock forward by d.
func (c *clock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// at parses an RFC 3339 time from the fixtures.
func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// pass returns the catalog's pass with this ID.
func pass(id, aos, los string, maxElevDeg float64) PassView {
	sat, rest, _ := strings.Cut(id, "/")
	station, _, _ := strings.Cut(rest, "/")
	return PassView{ID: id, Pass: Pass{SatID: sat, Station: station, AOS: at(aos), LOS: at(los), MaxElevDeg: maxElevDeg}}
}

// testService returns the booking service as Mirela reviews it on Friday morning, reading its time from now.
func testService(t *testing.T, now func() time.Time) *Service {
	t.Helper()
	store := NewMemStore()
	err := store.load([]Booking{
		{ID: "bk-4501", SatID: "furrow-a", Station: "SQT", AOS: at("2027-12-06T14:02:10Z"), LOS: at("2027-12-06T14:11:52Z")},
		{ID: "bk-4510", SatID: "cirrine-1", Station: "YST", AOS: at("2028-01-07T06:59:16Z"), LOS: at("2028-01-07T07:09:06Z")},
		{ID: "bk-4513", SatID: "brine-3", Station: "LAN", AOS: at("2028-01-07T10:12:19Z"), LOS: at("2028-01-07T10:19:56Z")},
		{ID: "bk-4514", SatID: "pipit-1", Station: "PEL", AOS: at("2028-01-07T04:00:03Z"), LOS: at("2028-01-07T04:09:14Z")},
		{ID: "bk-4516", SatID: "furrow-a", Station: "SQT", AOS: at("2028-01-07T14:26:42Z"), LOS: at("2028-01-07T14:36:45Z")},
	})
	if err != nil {
		t.Fatal(err)
	}
	return &Service{
		Stations: []Station{{Code: "SQT", Name: "Salar Quieto", LatDeg: -23.4, LonDeg: -104.6, MaskDeg: 3, Bands: []Band{"S", "X"}}},
		Catalog: []PassView{
			pass("furrow-a/SQT/16351", "2028-01-07T14:26:42Z", "2028-01-07T14:36:45Z", 66.1),
			pass("furrow-b/SQT/16349", "2028-01-07T15:14:45Z", "2028-01-07T15:23:31Z", 19.0),
			pass("cirrine-1/YST/11856", "2028-01-08T07:07:18Z", "2028-01-08T07:17:12Z", 37.4),
			pass("brine-3/LAN/09719", "2028-01-08T09:42:47Z", "2028-01-08T09:50:38Z", 45.7),
		},
		Guide: fstest.MapFS{"index.html": {Data: []byte("<h1>Booking passes with Apogee</h1>\n")}},
		Store: store,
		Customers: map[string]Customer{
			"apogee":      {ID: "apogee", Name: "Apogee", SatIDs: []string{"pipit-1", "pipit-2", "pipit-3"}, KeyHash: "8018b533030c734641dda97af10c84dd57357d0e37b7a0bdee7ed2a63c911407"},
			"furrowlight": {ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}, KeyHash: "95ff036b9ad2e59c26fae1be7dda78b12412d243fcbb6567bc8124e6f1c93fd3"},
			"brinewatch":  {ID: "brinewatch", Name: "Brinewatch", SatIDs: []string{"brine-3"}, KeyHash: "6f5748f2886bd91e86930360f00a38ae260f6de18ea20012f4e4ec718d483f85"},
			"cirrine":     {ID: "cirrine", Name: "Cirrine Atmospherics", SatIDs: []string{"cirrine-1"}, KeyHash: "d9ee2e5de1796bbeba223c0737933b3e60d53687d479f8a04c98c99df87b510d"},
		},
		TokenKey: []byte(tokenKey),
		Now:      now,
	}
}

// bugs are the four broken services Mirela planted, each named for what is wrong with it.
var bugs = []string{
	"fails open",
	"lists every customer's bookings",
	"never lets a token expire",
	"takes an API key as a token",
}

// plant returns s's routes with one of Mirela's bugs in them: the private routes go through a broken
// guard or a broken handler, and every other request goes to the real routes.
func plant(s *Service, bug string) http.Handler {
	guard := s.private
	listBookings := s.listBookings
	switch bug {
	case "fails open":
		guard = func(h http.HandlerFunc) http.Handler {
			checked := s.private(h)
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") == "" {
					h(w, r) // the bug: no credentials, so nothing to check
					return
				}
				checked.ServeHTTP(w, r)
			})
		}
	case "lists every customer's bookings":
		listBookings = func(w http.ResponseWriter, r *http.Request) {
			all, err := s.Store.List(r.Context()) // the bug: nobody asks owns
			if err != nil {
				writeError(w, r, err)
				return
			}
			writeJSON(w, http.StatusOK, all)
		}
	case "never lets a token expire":
		started := s.Now() // the bug: the guard's clock stopped when the service started
		guard = func(h http.HandlerFunc) http.Handler {
			return authenticate(s.Customers, s.TokenKey, func() time.Time { return started })(h)
		}
	case "takes an API key as a token":
		guard = func(h http.HandlerFunc) http.Handler {
			checked := s.private(h)
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				key, err := credentials(r, "Bearer")
				if c, ok := checkKey(s.Customers, key); err == nil && ok {
					h(w, r.WithContext(withCustomer(r.Context(), c))) // the bug: a key where a token belongs
					return
				}
				checked.ServeHTTP(w, r)
			})
		}
	}
	mux := http.NewServeMux()
	mux.Handle("GET /passes", guard(s.listPasses))
	mux.Handle("GET /passes/{sat}/{station}/{orbit}", guard(s.getPass))
	mux.Handle("POST /bookings", guard(s.createBooking))
	mux.Handle("GET /bookings", guard(listBookings))
	mux.Handle("GET /bookings/{id}", guard(s.getBooking))
	mux.Handle("DELETE /bookings/{id}", guard(s.cancelBooking))
	mux.Handle("/", s.routes())
	return mux
}

// serve starts the service, with bug planted in it unless bug is "", and returns its address and
// the function that moves its clock forward.
func serve(t *testing.T, bug string) (string, func(time.Duration)) {
	c := &clock{now: time.Date(2028, 1, 7, 9, 0, 0, 0, time.UTC)}
	s := testService(t, c.Now)
	var h http.Handler = s.routes()
	if bug != "" {
		h = plant(s, bug)
	}
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv.URL, c.Add
}

func TestReviewPassesOnTheService(t *testing.T) {
	base, advance := serve(t, "")
	review(t, base, advance)
}

// TestMutant runs the review against the planted service APOGEE_MUTANT names. In a normal run it
// does nothing: TestReviewCatchesEveryBug runs it in a fresh copy of this test binary.
func TestMutant(t *testing.T) {
	bug := os.Getenv("APOGEE_MUTANT")
	if bug == "" {
		return
	}
	if !slices.Contains(bugs, bug) {
		t.Fatalf("no planted service %q", bug)
	}
	base, advance := serve(t, bug)
	review(t, base, advance)
}

// TestReviewCatchesEveryBug runs the review against each planted service in a fresh copy of this
// test binary. Every one of those runs must fail.
func TestReviewCatchesEveryBug(t *testing.T) {
	for _, bug := range bugs {
		cmd := exec.Command(os.Args[0], "-test.run=^TestMutant$", "-test.timeout=3s")
		cmd.Env = append(os.Environ(), "APOGEE_MUTANT="+bug)
		out, err := cmd.CombinedOutput()
		var exit *exec.ExitError
		switch {
		case err == nil:
			t.Errorf("review against the service that %s = passed, want a failed check: none of your checks catches it", bug)
		case !errors.As(err, &exit):
			t.Fatalf("running the review against the service that %s: %v", bug, err)
		case !strings.Contains(string(out), "--- FAIL: TestMutant"):
			first, _, _ := strings.Cut(strings.TrimSpace(string(out)), "\n")
			t.Errorf("review against the service that %s = %s, want a failed check", bug, first)
		}
	}
}
