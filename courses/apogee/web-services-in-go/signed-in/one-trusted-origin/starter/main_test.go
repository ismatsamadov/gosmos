package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// planningPage is the origin of Furrowlight's planning page: the one origin the API trusts.
const planningPage = "https://plan.furrowlight.example"

// token is a token Furrowlight's server got from POST /tokens, valid until 18:00 UTC; the stand-in API doesn't check it.
const token = "furrowlight.1832090400.8e8cd6d144a78b2abd3e199c0ac7b992109fe90cd539f3faa80010852fbcbb71"

// api stands in for the booking service's routes. It refuses a request without a token, as the service does, and counts every request it sees.
func api(seen *int) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /passes", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, `{"error":"unauthorized","message":"missing credentials"}`, http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `[{"id":"furrow-a/SQT/16609"}]`+"\n")
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*seen++
		mux.ServeHTTP(w, r)
	})
}

// call sends one request from a page at origin, or from a program when origin is empty, through cors around the API.
func call(method, origin string, header map[string]string) (*http.Response, string, int) {
	seen := 0
	r := httptest.NewRequest(method, "/passes", nil)
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	for k, v := range header {
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	cors(planningPage)(api(&seen)).ServeHTTP(rec, r)
	res := rec.Result()
	body, _ := io.ReadAll(res.Body)
	return res, strings.TrimSpace(string(body)), seen
}

// variesByOrigin reports whether an answer's Vary header names Origin.
func variesByOrigin(res *http.Response) bool {
	for _, v := range res.Header.Values("Vary") {
		for _, name := range strings.Split(v, ",") {
			if strings.EqualFold(strings.TrimSpace(name), "Origin") {
				return true
			}
		}
	}
	return false
}

// preflight is what the browser asks before the planning page's script may send GET /passes with a token.
var preflight = map[string]string{
	"Access-Control-Request-Method":  "GET",
	"Access-Control-Request-Headers": "authorization",
}

func TestPreflight(t *testing.T) {
	res, _, seen := call(http.MethodOptions, planningPage, preflight)
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("preflight OPTIONS /passes from %s = %d, want 204", planningPage, res.StatusCode)
	}
	want := map[string]string{
		"Access-Control-Allow-Origin":  planningPage,
		"Access-Control-Allow-Methods": "GET, POST, DELETE",
		"Access-Control-Allow-Headers": "Authorization, Content-Type, Idempotency-Key",
	}
	for _, name := range []string{"Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers"} {
		if got := res.Header.Get(name); got != want[name] {
			t.Errorf("preflight from %s: %s = %q, want %q", planningPage, name, got, want[name])
		}
	}
	if !variesByOrigin(res) {
		t.Errorf("preflight from %s: Vary = %q, want it to name Origin", planningPage, res.Header.Values("Vary"))
	}
	if seen != 0 {
		t.Errorf("the preflight reached the API %d time(s), want 0: it carries no token, so the API would refuse it", seen)
	}
}

func TestThePlanningPageReads(t *testing.T) {
	res, body, _ := call(http.MethodGet, planningPage, map[string]string{"Authorization": "Bearer " + token})
	if res.StatusCode != http.StatusOK || body != `[{"id":"furrow-a/SQT/16609"}]` {
		t.Fatalf("GET /passes from %s = %d %s, want the API's 200 and its passes", planningPage, res.StatusCode, body)
	}
	if got := res.Header.Get("Access-Control-Allow-Origin"); got != planningPage {
		t.Errorf("GET /passes from %s: Access-Control-Allow-Origin = %q, want %q", planningPage, got, planningPage)
	}
	if !variesByOrigin(res) {
		t.Errorf("GET /passes from %s: Vary = %q, want it to name Origin", planningPage, res.Header.Values("Vary"))
	}
}

func TestOtherOrigins(t *testing.T) {
	for _, origin := range []string{
		"https://elsewhere.example",
		"http://plan.furrowlight.example",
		"https://plan.furrowlight.example.elsewhere.example",
		"https://www.furrowlight.example",
	} {
		res, _, _ := call(http.MethodOptions, origin, preflight)
		if got := res.Header.Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("preflight from %s: Access-Control-Allow-Origin = %q, want none", origin, got)
		}
		res, body, _ := call(http.MethodGet, origin, map[string]string{"Authorization": "Bearer " + token})
		if res.StatusCode != http.StatusOK || body != `[{"id":"furrow-a/SQT/16609"}]` {
			t.Errorf("GET /passes from %s = %d %s, want the API's 200 and its passes: the browser, not the service, hides them", origin, res.StatusCode, body)
		}
		if got := res.Header.Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("GET /passes from %s: Access-Control-Allow-Origin = %q, want none", origin, got)
		}
		if !variesByOrigin(res) {
			t.Errorf("GET /passes from %s: Vary = %q, want it to name Origin", origin, res.Header.Values("Vary"))
		}
	}
}

func TestProgramsAreUntouched(t *testing.T) {
	res, body, _ := call(http.MethodGet, "", map[string]string{"Authorization": "Bearer " + token})
	if res.StatusCode != http.StatusOK || body != `[{"id":"furrow-a/SQT/16609"}]` {
		t.Fatalf("GET /passes from a program = %d %s, want the API's 200 and its passes", res.StatusCode, body)
	}
	if got := res.Header.Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("GET /passes from a program: Access-Control-Allow-Origin = %q, want none", got)
	}
	if !variesByOrigin(res) {
		t.Errorf("GET /passes from a program: Vary = %q, want it to name Origin", res.Header.Values("Vary"))
	}
	res, _, _ = call(http.MethodGet, "", nil)
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("GET /passes with no token = %d, want the API's own 401", res.StatusCode)
	}
}
