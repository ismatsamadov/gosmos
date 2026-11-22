package main

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"testing"
	"testing/fstest"
	"time"
)

// monday is when the full rehearsal runs: 10:00 UTC on Monday 31 January.
var monday = time.Date(2028, 1, 31, 10, 0, 0, 0, time.UTC)

// rehearsal builds the whole service as Monday's rehearsal runs it, on a clock that stays at 10:00.
func rehearsal(t *testing.T) *App {
	t.Helper()
	app, err := newApp(Config{
		Dir:       t.TempDir(),
		Stations:  testStations(),
		Catalog:   testCatalog(),
		Guide:     fstest.MapFS{"index.html": {Data: []byte("<h1>Booking passes with Apogee</h1>\n")}},
		Customers: testCustomers(),
		TokenKey:  []byte(testTokenKey),
		Now:       func() time.Time { return monday },
		Rand:      rand.New(rand.NewPCG(1, 2)),
	})
	if err != nil {
		t.Fatalf("newApp = %v, want nil", err)
	}
	return app
}

// serveApp serves h on a test server on 127.0.0.1 and returns its address.
func serveApp(t *testing.T, h http.Handler) string {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv.URL
}

// without removes one header from every request before h sees it.
func without(name string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del(name)
		h.ServeHTTP(w, r)
	})
}

// plainErrors answers every error in plain text, as the service did before its errors were JSON.
func plainErrors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		for name, values := range rec.Header() {
			w.Header()[name] = values
		}
		if rec.Code < 400 {
			w.WriteHeader(rec.Code)
			w.Write(rec.Body.Bytes())
			return
		}
		var body ErrorBody
		json.Unmarshal(rec.Body.Bytes(), &body)
		http.Error(w, body.Message, rec.Code)
	})
}

// alwaysLater is a clock a minute later at every reading: a bucket on it is full again at every request.
func alwaysLater() func() time.Time {
	var readings atomic.Int64
	return func() time.Time { return monday.Add(time.Duration(readings.Add(1)) * time.Minute) }
}

// mutants are Nkem's planted services, each with one of the season's bugs.
var mutants = []struct {
	bug   string
	plant func(app *App)
}{
	{"ignores idempotency keys, so a retried booking gets 409", func(app *App) {
		app.Handler = without("Idempotency-Key", app.Handler)
	}},
	{"shows every customer's bookings", func(app *App) {
		brinewatch := app.Service.Customers["brinewatch"]
		brinewatch.SatIDs = []string{"pipit-1", "pipit-2", "pipit-3", "furrow-a", "furrow-b", "brine-3", "cirrine-1"}
		app.Service.Customers["brinewatch"] = brinewatch
	}},
	{"never limits", func(app *App) {
		app.Service.Limiter.now = alwaysLater()
	}},
	{"answers errors as plain text", func(app *App) {
		app.Handler = plainErrors(app.Handler)
	}},
}

func TestRealServicePasses(t *testing.T) {
	checkService(t, serveApp(t, rehearsal(t).Handler))
}

// TestMutant runs checkService against the planted service APOGEE_MUTANT names. In a normal run it does
// nothing: TestCheckCatchesMutants runs it in a fresh copy of this test binary.
func TestMutant(t *testing.T) {
	bug := os.Getenv("APOGEE_MUTANT")
	if bug == "" {
		return
	}
	for _, m := range mutants {
		if m.bug == bug {
			app := rehearsal(t)
			m.plant(app)
			checkService(t, serveApp(t, app.Handler))
			return
		}
	}
	t.Fatalf("no mutant %q", bug)
}

// TestCheckCatchesMutants runs checkService against each planted service in a fresh copy of this test
// binary. Every run must fail.
func TestCheckCatchesMutants(t *testing.T) {
	for _, m := range mutants {
		cmd := exec.Command(os.Args[0], "-test.run=^TestMutant$", "-test.count=1", "-test.timeout=5s")
		cmd.Env = append(os.Environ(), "APOGEE_MUTANT="+m.bug)
		out, err := cmd.CombinedOutput()
		var exit *exec.ExitError
		switch {
		case err == nil:
			t.Errorf("the service that %s passed checkService, want checkService to fail on it", m.bug)
		case !errors.As(err, &exit):
			t.Fatalf("running the service that %s: %v", m.bug, err)
		case !strings.Contains(string(out), "--- FAIL: TestMutant"):
			t.Errorf("the service that %s didn't fail a check, it crashed:\n%s", m.bug, out)
		}
	}
}
