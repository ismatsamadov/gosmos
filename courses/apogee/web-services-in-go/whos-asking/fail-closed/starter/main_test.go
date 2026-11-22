package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

// furrowlightKey is Furrowlight's test API key.
const furrowlightKey = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"

// tokenKey is the test token key: 32 bytes, like the service's real key, which never appears in a lesson.
const tokenKey = "bookingd test token key 32 bytes"

// now is the service's clock in these tests: Thursday 6 January, 10:30 UTC.
var now = time.Date(2028, 1, 6, 10, 30, 0, 0, time.UTC)

// eleven is when this morning's tokens expire.
var eleven = time.Date(2028, 1, 6, 11, 0, 0, 0, time.UTC)

// customers returns the customer table as the service keeps it: a hash of each key, never the key.
func customers() map[string]Customer {
	return map[string]Customer{
		"apogee":      {ID: "apogee", Name: "Apogee", SatIDs: []string{"pipit-1", "pipit-2", "pipit-3"}, KeyHash: "8018b533030c734641dda97af10c84dd57357d0e37b7a0bdee7ed2a63c911407"},
		"furrowlight": {ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}, KeyHash: "95ff036b9ad2e59c26fae1be7dda78b12412d243fcbb6567bc8124e6f1c93fd3"},
		"brinewatch":  {ID: "brinewatch", Name: "Brinewatch", SatIDs: []string{"brine-3"}, KeyHash: "6f5748f2886bd91e86930360f00a38ae260f6de18ea20012f4e4ec718d483f85"},
		"cirrine":     {ID: "cirrine", Name: "Cirrine Atmospherics", SatIDs: []string{"cirrine-1"}, KeyHash: "d9ee2e5de1796bbeba223c0737933b3e60d53687d479f8a04c98c99df87b510d"},
	}
}

// mint signs a token the way the service does, worked out here so the test needn't trust signToken.
func mint(customer string, expires time.Time, key string) string {
	claims := customer + "." + strconv.FormatInt(expires.Unix(), 10)
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(claims))
	return claims + "." + hex.EncodeToString(mac.Sum(nil))
}

// at parses an RFC 3339 time from the fixtures.
func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// testService returns the booking service as it stands this morning, with one booking of Cirrine's.
func testService(t *testing.T) *Service {
	t.Helper()
	store := NewMemStore()
	err := store.load([]Booking{
		{ID: "bk-4510", SatID: "cirrine-1", Station: "YST", AOS: at("2028-01-07T06:59:16Z"), LOS: at("2028-01-07T07:09:06Z")},
	})
	if err != nil {
		t.Fatal(err)
	}
	return &Service{
		Stations: []Station{
			{Code: "YST", Name: "Ystmark", LatDeg: 77.8, LonDeg: 2.0, MaskDeg: 5, Bands: []Band{"S", "X"}},
			{Code: "SQT", Name: "Salar Quieto", LatDeg: -23.4, LonDeg: -104.6, MaskDeg: 3, Bands: []Band{"S", "X"}},
		},
		Catalog: []PassView{
			{ID: "cirrine-1/YST/11841", Pass: Pass{SatID: "cirrine-1", Station: "YST", AOS: at("2028-01-07T06:59:16Z"), LOS: at("2028-01-07T07:09:06Z"), MaxElevDeg: 35.3}},
			{ID: "furrow-a/SQT/16351", Pass: Pass{SatID: "furrow-a", Station: "SQT", AOS: at("2028-01-07T14:26:42Z"), LOS: at("2028-01-07T14:36:45Z"), MaxElevDeg: 66.1}},
		},
		Guide:     fstest.MapFS{"index.html": {Data: []byte("<h1>Booking passes with Apogee</h1>\n")}},
		Store:     store,
		Customers: customers(),
		TokenKey:  []byte(tokenKey),
		Now:       func() time.Time { return now },
	}
}

// request returns a request with auth as its Authorization header; "" sends none.
func request(method, path, auth string) *http.Request {
	r := httptest.NewRequest(method, path, nil)
	if auth != "" {
		r.Header.Set("Authorization", auth)
	}
	return r
}

func TestAuthenticateRefuses(t *testing.T) {
	sig := mint("furrowlight", eleven, tokenKey)[len("furrowlight.1830769200."):]
	tests := []struct {
		name string
		auth string // "" sends no Authorization header
		want string // the message of the 401
	}{
		{"no Authorization header", "", "missing credentials"},
		{"an API key where a token belongs", "Key " + furrowlightKey, "missing credentials"},
		{"the API key sent as a Bearer token", "Bearer " + furrowlightKey, "invalid token"},
		{"Furrowlight's token edited to say cirrine", "Bearer cirrine.1830769200." + sig, "invalid token"},
		{"a token for a customer the service doesn't know", "Bearer " + mint("kessick", eleven, tokenKey), "invalid token"},
		{"a token that expired at 10:30", "Bearer " + mint("furrowlight", now, tokenKey), "token expired"},
	}
	guard := authenticate(customers(), []byte(tokenKey), func() time.Time { return now })
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.Write([]byte("every booking there is\n"))
			})
			rec := httptest.NewRecorder()
			guard(next).ServeHTTP(rec, request(http.MethodGet, "/bookings", tt.auth))
			if calls != 0 {
				t.Errorf("GET /bookings with %s: next ran %d times, want 0: the guard must return after a 401", tt.name, calls)
			}
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("GET /bookings with %s = %d, want 401", tt.name, rec.Code)
			}
			if got := rec.Header().Get("WWW-Authenticate"); got != "Bearer" {
				t.Errorf("GET /bookings with %s WWW-Authenticate = %q, want %q", tt.name, got, "Bearer")
			}
			var got ErrorBody
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("GET /bookings with %s body: %v, want an ErrorBody", tt.name, err)
			}
			if want := (ErrorBody{Code: "unauthorized", Message: tt.want}); got != want {
				t.Errorf("GET /bookings with %s = %+v, want %+v", tt.name, got, want)
			}
		})
	}
}

func TestAuthenticatePassesTheCustomerOn(t *testing.T) {
	guard := authenticate(customers(), []byte(tokenKey), func() time.Time { return now })
	calls := 0
	var seen Customer
	var found bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		seen, found = customerFrom(r.Context())
	})
	rec := httptest.NewRecorder()
	guard(next).ServeHTTP(rec, request(http.MethodGet, "/bookings", "Bearer "+mint("furrowlight", eleven, tokenKey)))
	if calls != 1 {
		t.Fatalf("GET /bookings with Furrowlight's token: next ran %d times, want 1", calls)
	}
	if !found || seen.ID != "furrowlight" || seen.Name != "Furrowlight" {
		t.Errorf("customerFrom(the context next got) = %q %q, %v, want %q %q, true", seen.ID, seen.Name, found, "furrowlight", "Furrowlight")
	}
	if rec.Code != http.StatusOK || rec.Header().Get("WWW-Authenticate") != "" {
		t.Errorf("GET /bookings with Furrowlight's token = %d, WWW-Authenticate %q, want 200 and none", rec.Code, rec.Header().Get("WWW-Authenticate"))
	}
}

func TestOnlyPrivateRoutesNeedAToken(t *testing.T) {
	h := testService(t).routes()
	tests := []struct {
		route string
		auth  string
		want  int
	}{
		{"GET /health", "", http.StatusOK},
		{"GET /stations", "", http.StatusOK},
		{"GET /stations/SQT", "", http.StatusOK},
		{"GET /guide/", "", http.StatusOK},
		{"GET /passes", "", http.StatusUnauthorized},
		{"GET /passes/furrow-a/SQT/16351", "", http.StatusUnauthorized},
		{"POST /bookings", "", http.StatusUnauthorized},
		{"GET /bookings", "", http.StatusUnauthorized},
		{"GET /bookings/bk-4510", "", http.StatusUnauthorized},
		{"DELETE /bookings/bk-4510", "", http.StatusUnauthorized},
		{"GET /bookings", "Bearer " + mint("cirrine", eleven, tokenKey), http.StatusOK},
		{"GET /bookings/bk-4510", "Bearer " + mint("cirrine", eleven, tokenKey), http.StatusOK},
	}
	for _, tt := range tests {
		method, path, _ := strings.Cut(tt.route, " ")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, request(method, path, tt.auth))
		with := "with no token"
		if tt.auth != "" {
			with = "with Cirrine's token"
		}
		if rec.Code != tt.want {
			t.Errorf("%s %s = %d, want %d", tt.route, with, rec.Code, tt.want)
		}
	}
}
