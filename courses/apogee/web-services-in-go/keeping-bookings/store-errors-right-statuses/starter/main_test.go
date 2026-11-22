package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// january holds Furrowlight's first pass of January over Salar Quieto.
var january = []PassView{
	{ID: "furrow-a/SQT/16290", Pass: Pass{SatID: "furrow-a", Station: "SQT",
		AOS: time.Date(2028, 1, 3, 14, 11, 55, 0, time.UTC), LOS: time.Date(2028, 1, 3, 14, 21, 59, 0, time.UTC), MaxElevDeg: 70.5}},
}

// thursday is the service's clock: Thursday 16 December 2027.
func thursday() time.Time { return time.Date(2027, 12, 16, 10, 0, 0, 0, time.UTC) }

const (
	furrowA   = `{"sat_id":"furrow-a","station":"SQT","aos":"2028-01-03T14:11:55Z","los":"2028-01-03T14:21:59Z"}`
	backwards = `{"sat_id":"furrow-a","station":"SQT","aos":"2028-01-03T14:21:59Z","los":"2028-01-03T14:11:55Z"}`
)

// errFull is what a full disk says; no caller may ever see it.
var errFull = errors.New("write bookings.json: no space left on device")

// fullDisk is a store whose disk is full: it keeps nothing.
type fullDisk struct{ *MemStore }

// Create returns errFull.
func (fullDisk) Create(ctx context.Context, b Booking) (Booking, error) { return Booking{}, errFull }

// checkError checks that an answer is an ErrorBody with this status, code and message.
func checkError(t *testing.T, call string, status int, header http.Header, body []byte, want ErrorBody, wantStatus int) {
	t.Helper()
	if status != wantStatus {
		t.Errorf("%s = %d, want %d", call, status, wantStatus)
	}
	if got := header.Get("Content-Type"); got != "application/json" {
		t.Errorf("%s: Content-Type = %q, want %q", call, got, "application/json")
	}
	var got ErrorBody
	if err := json.Unmarshal(body, &got); err != nil {
		t.Errorf("%s: body %q doesn't decode as an ErrorBody: %v", call, body, err)
		return
	}
	got.RequestID = "" // the IDs are chapter 6's business
	if got != want {
		t.Errorf("%s: body = %+v, want %+v", call, got, want)
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		want   ErrorBody
	}{
		{"a booking that doesn't exist", fmt.Errorf("booking %s: %w", "bk-5499", ErrNotFound),
			http.StatusNotFound, ErrorBody{Code: "not_found", Message: "booking bk-5499: not found"}},
		{"a pass already booked", fmt.Errorf("furrow-a over SQT at 2028-01-03T14:11:55Z: %w", ErrConflict),
			http.StatusConflict, ErrorBody{Code: "conflict", Message: "furrow-a over SQT at 2028-01-03T14:11:55Z: already booked"}},
		{"a full disk", errFull,
			http.StatusInternalServerError, ErrorBody{Code: "internal", Message: "internal error"}},
		{"a broken booking rule", &APIError{Status: http.StatusUnprocessableEntity, Code: "invalid_booking", Message: "los must be after aos"},
			http.StatusUnprocessableEntity, ErrorBody{Code: "invalid_booking", Message: "los must be after aos"}},
		{"an API error beside ErrNotFound", errors.Join(&APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: `unknown band "x"`}, ErrNotFound),
			http.StatusBadRequest, ErrorBody{Code: "bad_request", Message: `unknown band "x"`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeError(rec, httptest.NewRequest(http.MethodGet, "/bookings/bk-5499", nil), tt.err)
			checkError(t, fmt.Sprintf("writeError(%q)", tt.err.Error()), rec.Code, rec.Header(), rec.Body.Bytes(), tt.want, tt.status)
		})
	}
}

// send sends one request through h and returns its status, headers and body.
func send(h http.Handler, method, path, body string) (int, http.Header, []byte) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Header(), rec.Body.Bytes()
}

func TestStoreErrorsOverHTTP(t *testing.T) {
	s := &Service{Catalog: january, Store: NewMemStore(), Now: thursday}
	h := s.handler(log.New(io.Discard, "", 0))
	notFound := ErrorBody{Code: "not_found", Message: "booking bk-5499: not found"}

	status, header, body := send(h, http.MethodGet, "/bookings/bk-5499", "")
	checkError(t, "GET /bookings/bk-5499", status, header, body, notFound, http.StatusNotFound)
	status, header, body = send(h, http.MethodDelete, "/bookings/bk-5499", "")
	checkError(t, "DELETE /bookings/bk-5499", status, header, body, notFound, http.StatusNotFound)

	if status, _, body := send(h, http.MethodPost, "/bookings", furrowA); status != http.StatusCreated {
		t.Fatalf("POST /bookings (Furrow-A on 3 January) = %d %s, want %d", status, body, http.StatusCreated)
	}
	status, header, body = send(h, http.MethodPost, "/bookings", furrowA)
	checkError(t, "POST /bookings (the same pass again)", status, header, body,
		ErrorBody{Code: "conflict", Message: "furrow-a over SQT at 2028-01-03T14:11:55Z: already booked"}, http.StatusConflict)
	status, header, body = send(h, http.MethodPost, "/bookings", backwards)
	checkError(t, "POST /bookings (LOS before AOS)", status, header, body,
		ErrorBody{Code: "invalid_booking", Message: "los must be after aos"}, http.StatusUnprocessableEntity)

	full := &Service{Catalog: january, Store: fullDisk{NewMemStore()}, Now: thursday}
	status, header, body = send(full.handler(log.New(io.Discard, "", 0)), http.MethodPost, "/bookings", furrowA)
	checkError(t, "POST /bookings (a full disk)", status, header, body, ErrorBody{Code: "internal", Message: "internal error"}, http.StatusInternalServerError)
}

func TestHandlersLeaveStoreErrorsToWriteError(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "bookings.go", nil, 0)
	if err != nil {
		t.Fatalf("reading bookings.go: %v", err)
	}
	ast.Inspect(file, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && (id.Name == "ErrNotFound" || id.Name == "ErrConflict") {
			t.Errorf("bookings.go:%d uses %s; want writeError to be the one place that turns store errors into statuses", fset.Position(id.Pos()).Line, id.Name)
		}
		return true
	})
}
