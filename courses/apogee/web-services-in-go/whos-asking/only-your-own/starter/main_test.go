package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

// tokenKey is the test token key: 32 bytes, like the service's real key, which never appears in a lesson.
const tokenKey = "bookingd test token key 32 bytes"

// now is the service's clock in these tests: Thursday 6 January, 10:00 UTC.
var now = time.Date(2028, 1, 6, 10, 0, 0, 0, time.UTC)

// customers returns the customer table as the service keeps it: a hash of each key, never the key.
func customers() map[string]Customer {
	return map[string]Customer{
		"apogee":      {ID: "apogee", Name: "Apogee", SatIDs: []string{"pipit-1", "pipit-2", "pipit-3"}, KeyHash: "8018b533030c734641dda97af10c84dd57357d0e37b7a0bdee7ed2a63c911407"},
		"furrowlight": {ID: "furrowlight", Name: "Furrowlight", SatIDs: []string{"furrow-a", "furrow-b"}, KeyHash: "95ff036b9ad2e59c26fae1be7dda78b12412d243fcbb6567bc8124e6f1c93fd3"},
		"brinewatch":  {ID: "brinewatch", Name: "Brinewatch", SatIDs: []string{"brine-3"}, KeyHash: "6f5748f2886bd91e86930360f00a38ae260f6de18ea20012f4e4ec718d483f85"},
		"cirrine":     {ID: "cirrine", Name: "Cirrine Atmospherics", SatIDs: []string{"cirrine-1"}, KeyHash: "d9ee2e5de1796bbeba223c0737933b3e60d53687d479f8a04c98c99df87b510d"},
	}
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

// testRoutes returns the routes of the booking service as it stands this Thursday: a week's catalog and five bookings.
func testRoutes(t *testing.T) http.Handler {
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
	s := &Service{
		Stations: []Station{{Code: "SQT", Name: "Salar Quieto", LatDeg: -23.4, LonDeg: -104.6, MaskDeg: 3, Bands: []Band{"S", "X"}}},
		Catalog: []PassView{
			pass("pipit-1/PEL/08639", "2028-01-07T04:00:03Z", "2028-01-07T04:09:14Z", 43.6),
			pass("cirrine-1/YST/11841", "2028-01-07T06:59:16Z", "2028-01-07T07:09:06Z", 35.3),
			pass("brine-3/LAN/09704", "2028-01-07T10:12:19Z", "2028-01-07T10:19:56Z", 35.7),
			pass("furrow-a/SQT/16351", "2028-01-07T14:26:42Z", "2028-01-07T14:36:45Z", 66.1),
			pass("furrow-b/SQT/16349", "2028-01-07T15:14:45Z", "2028-01-07T15:23:31Z", 19.0),
			pass("cirrine-1/YST/11856", "2028-01-08T07:07:18Z", "2028-01-08T07:17:12Z", 37.4),
			pass("brine-3/LAN/09719", "2028-01-08T09:42:47Z", "2028-01-08T09:50:38Z", 45.7),
		},
		Guide:     fstest.MapFS{"index.html": {Data: []byte("<h1>Booking passes with Apogee</h1>\n")}},
		Store:     store,
		Customers: customers(),
		TokenKey:  []byte(tokenKey),
		Now:       func() time.Time { return now },
	}
	return s.routes()
}

// tokenFor signs a token for customer, valid for this morning, the way the service does.
func tokenFor(customer string) string {
	claims := customer + "." + strconv.FormatInt(now.Add(time.Hour).Unix(), 10)
	mac := hmac.New(sha256.New, []byte(tokenKey))
	mac.Write([]byte(claims))
	return claims + "." + hex.EncodeToString(mac.Sum(nil))
}

// ask sends route, such as "GET /bookings", to h as customer, with body as JSON when it isn't "".
func ask(h http.Handler, customer, route, body string) *httptest.ResponseRecorder {
	method, path, _ := strings.Cut(route, " ")
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, r)
	req.Header.Set("Authorization", "Bearer "+tokenFor(customer))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// ids returns the IDs of the list a 200 answer holds, or what came back instead.
func ids(rec *httptest.ResponseRecorder) string {
	var list []struct {
		ID string `json:"id"`
	}
	if rec.Code != http.StatusOK || json.Unmarshal(rec.Body.Bytes(), &list) != nil || list == nil {
		return strconv.Itoa(rec.Code) + " " + string(bytes.TrimSpace(rec.Body.Bytes()))
	}
	out := []string{}
	for _, item := range list {
		out = append(out, item.ID)
	}
	return "[" + strings.Join(out, " ") + "]"
}

// answer returns an answer's status and its ErrorBody, for the messages of a failed check.
func answer(rec *httptest.ResponseRecorder) string {
	var body ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Code == "" {
		return strconv.Itoa(rec.Code) + " " + string(bytes.TrimSpace(rec.Body.Bytes()))
	}
	return strconv.Itoa(rec.Code) + " " + body.Code + " " + strconv.Quote(body.Message)
}

func TestOwns(t *testing.T) {
	c := customers()
	tests := []struct {
		name     string
		customer Customer
		satID    string
		want     bool
	}{
		{"Furrowlight and Furrow-A", c["furrowlight"], "furrow-a", true},
		{"Furrowlight and Furrow-B", c["furrowlight"], "furrow-b", true},
		{"Furrowlight and Brine-3", c["furrowlight"], "brine-3", false},
		{"Brinewatch and Brine-3", c["brinewatch"], "brine-3", true},
		{"Brinewatch and Cirrine-1", c["brinewatch"], "cirrine-1", false},
		{"Apogee and Pipit-2", c["apogee"], "pipit-2", true},
		{"Apogee and Furrow-A", c["apogee"], "furrow-a", false},
		{"Furrowlight and half a satellite ID", c["furrowlight"], "furrow", false},
		{"a customer with no satellites", Customer{}, "pipit-1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := owns(tt.customer, tt.satID); got != tt.want {
				t.Errorf("owns(%q, %q) = %v, want %v", tt.customer.ID, tt.satID, got, tt.want)
			}
		})
	}
}

func TestListsHoldOnlyYourOwn(t *testing.T) {
	h := testRoutes(t)
	tests := []struct {
		customer, route, want string
	}{
		{"brinewatch", "GET /passes", "[brine-3/LAN/09704 brine-3/LAN/09719]"},
		{"brinewatch", "GET /passes?sat=cirrine-1", "[]"},
		{"brinewatch", "GET /bookings", "[bk-4513]"},
		{"cirrine", "GET /passes", "[cirrine-1/YST/11841 cirrine-1/YST/11856]"},
		{"cirrine", "GET /bookings", "[bk-4510]"},
		{"furrowlight", "GET /passes", "[furrow-a/SQT/16351 furrow-b/SQT/16349]"},
		{"furrowlight", "GET /bookings", "[bk-4501 bk-4516]"},
		{"apogee", "GET /passes?station=PEL", "[pipit-1/PEL/08639]"},
		{"apogee", "GET /bookings", "[bk-4514]"},
	}
	for _, tt := range tests {
		if got := ids(ask(h, tt.customer, tt.route, "")); got != tt.want {
			t.Errorf("%s as %s = %s, want %s", tt.route, tt.customer, got, tt.want)
		}
	}
}

func TestAnotherCustomersIsMissing(t *testing.T) {
	h := testRoutes(t)
	tests := []struct {
		customer, route, want string
	}{
		{"brinewatch", "GET /bookings/bk-4599", `404 not_found "booking bk-4599: not found"`},
		{"brinewatch", "GET /bookings/bk-4510", `404 not_found "booking bk-4510: not found"`},
		{"brinewatch", "DELETE /bookings/bk-4510", `404 not_found "booking bk-4510: not found"`},
		{"brinewatch", "GET /passes/cirrine-1/YST/11841", `404 not_found "no pass cirrine-1/YST/11841"`},
	}
	for _, tt := range tests {
		if got := answer(ask(h, tt.customer, tt.route, "")); got != tt.want {
			t.Errorf("%s as %s = %s, want %s", tt.route, tt.customer, got, tt.want)
		}
	}
	for _, tt := range []struct{ customer, route string }{
		{"cirrine", "GET /bookings/bk-4510"},
		{"cirrine", "GET /passes/cirrine-1/YST/11841"},
		{"brinewatch", "GET /bookings/bk-4513"},
	} {
		if got := ask(h, tt.customer, tt.route, ""); got.Code != http.StatusOK {
			t.Errorf("%s as %s = %s, want 200: a customer still sees its own", tt.route, tt.customer, answer(got))
		}
	}
	if got := ask(h, "cirrine", "DELETE /bookings/bk-4510", ""); got.Code != http.StatusNoContent {
		t.Errorf("DELETE /bookings/bk-4510 as cirrine = %s, want 204", answer(got))
	}
}

func TestBookingAnotherCustomersSatellite(t *testing.T) {
	h := testRoutes(t)
	const cirrinePass = `{"sat_id":"cirrine-1","station":"YST","aos":"2028-01-08T07:07:18Z","los":"2028-01-08T07:17:12Z"}`
	const brinePass = `{"sat_id":"brine-3","station":"LAN","aos":"2028-01-08T09:42:47Z","los":"2028-01-08T09:50:38Z"}`
	const forbidden = `403 forbidden "cirrine-1 is not one of your satellites"`
	for _, tt := range []struct{ name, body, want string }{
		{"Cirrine-1's pass on 8 January", cirrinePass, forbidden},
		{"a Cirrine-1 pass the catalog doesn't have, answered like a real one", `{"sat_id":"cirrine-1","station":"YST","aos":"2028-01-08T07:30:00Z","los":"2028-01-08T07:40:00Z"}`, forbidden},
		{"a pass with no sat_id, which checkShape refuses first", `{"station":"LAN","aos":"2028-01-08T09:42:47Z","los":"2028-01-08T09:50:38Z"}`, `422 invalid_booking "sat_id is required"`},
	} {
		if got := answer(ask(h, "brinewatch", "POST /bookings", tt.body)); got != tt.want {
			t.Errorf("POST /bookings of %s, as brinewatch = %s, want %s", tt.name, got, tt.want)
		}
	}
	if got := ids(ask(h, "cirrine", "GET /bookings", "")); got != "[bk-4510]" {
		t.Errorf("GET /bookings as cirrine after Brinewatch's refused bookings = %s, want [bk-4510]: nothing stored", got)
	}
	for _, tt := range []struct{ customer, body, location string }{
		{"cirrine", cirrinePass, "/bookings/bk-4517"},
		{"brinewatch", brinePass, "/bookings/bk-4518"},
	} {
		rec := ask(h, tt.customer, "POST /bookings", tt.body)
		if rec.Code != http.StatusCreated || rec.Header().Get("Location") != tt.location {
			t.Errorf("POST /bookings of its own pass as %s = %s, Location %q, want 201, Location %q", tt.customer, answer(rec), rec.Header().Get("Location"), tt.location)
		}
	}
	if got := ids(ask(h, "cirrine", "GET /bookings", "")); got != "[bk-4510 bk-4517]" {
		t.Errorf("GET /bookings as cirrine = %s, want [bk-4510 bk-4517]", got)
	}
}
