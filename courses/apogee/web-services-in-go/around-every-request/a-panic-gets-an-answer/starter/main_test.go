package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mondaysGetPass fails as getPass did on Monday: it reads five digits of an orbit number that may have four.
func mondaysGetPass(w http.ResponseWriter, r *http.Request) {
	orbit := r.PathValue("orbit")[:5]
	writeJSON(w, http.StatusOK, map[string]string{"orbit": orbit})
}

// passes routes pass IDs to Monday's getPass, and anything under /cut to a handler that cuts the connection on purpose.
func passes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /passes/{sat}/{station}/{orbit}", mondaysGetPass)
	mux.HandleFunc("GET /cut", func(w http.ResponseWriter, r *http.Request) {
		panic(http.ErrAbortHandler)
	})
	return mux
}

// serve sends req to h and returns what h answered, or the value h panicked with.
func serve(h http.Handler, req *http.Request) (rec *httptest.ResponseRecorder, panicked any) {
	rec = httptest.NewRecorder()
	defer func() { panicked = recover() }()
	h.ServeHTTP(rec, req)
	return rec, nil
}

func TestRecoverPanics(t *testing.T) {
	var buf bytes.Buffer
	h := requestID(recoverPanics(log.New(&buf, "", 0))(passes()))
	for _, path := range []string{"/passes/cirrine-1/YST/11364", "/passes/furrow-a/SQT/15864"} {
		rec, panicked := serve(h, httptest.NewRequest(http.MethodGet, path, nil))
		if panicked != nil || rec.Code != http.StatusOK {
			t.Errorf("GET %s = %d, panic %v, want 200 and no panic", path, rec.Code, panicked)
		}
	}
	if buf.Len() != 0 {
		t.Errorf("after two requests that didn't panic, the log = %q, want nothing", buf.String())
	}

	const path = "/passes/cirrine-1/YST/1136"
	rec, panicked := serve(h, httptest.NewRequest(http.MethodGet, path, nil))
	if panicked != nil {
		t.Fatalf("GET %s panicked through recoverPanics with %v, want a 500 answer", path, panicked)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("GET %s = %d, want %d", path, rec.Code, http.StatusInternalServerError)
	}
	var body ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("GET %s body %q doesn't decode as an ErrorBody: %v", path, rec.Body.String(), err)
	}
	want := ErrorBody{Code: "internal", Message: "internal error", RequestID: "req-000003"}
	if body != want {
		t.Errorf("GET %s body = %+v, want %+v", path, body, want)
	}
	wantLog := "req-000003 panic: runtime error: slice bounds out of range [:5] with length 4\n"
	if buf.String() != wantLog {
		t.Errorf("log = %q, want %q", buf.String(), wantLog)
	}
}

func TestRecoverPanicsLetsAbortThrough(t *testing.T) {
	var buf bytes.Buffer
	h := recoverPanics(log.New(&buf, "", 0))(passes())
	rec, panicked := serve(h, httptest.NewRequest(http.MethodGet, "/cut", nil))
	if panicked != http.ErrAbortHandler {
		t.Errorf("GET /cut panicked with %v, want http.ErrAbortHandler to go on through recoverPanics", panicked)
	}
	if rec.Body.Len() != 0 || buf.Len() != 0 {
		t.Errorf("GET /cut wrote %q and logged %q, want neither", rec.Body.String(), buf.String())
	}
}

func TestRecoverPanicsOverTheWire(t *testing.T) {
	srv := httptest.NewUnstartedServer(requestID(recoverPanics(log.New(io.Discard, "", 0))(passes())))
	srv.Config.ErrorLog = log.New(io.Discard, "", 0) // the server's own report of a panic
	srv.Start()
	defer srv.Close()
	for _, tc := range []struct {
		path   string
		status int
	}{
		{"/passes/cirrine-1/YST/1136", http.StatusInternalServerError},
		{"/passes/cirrine-1/YST/11364", http.StatusOK},
	} {
		res, err := srv.Client().Get(srv.URL + tc.path)
		if err != nil {
			t.Errorf("GET %s error = %v, want an answer with status %d", tc.path, err, tc.status)
			continue
		}
		res.Body.Close()
		if res.StatusCode != tc.status {
			t.Errorf("GET %s = %d, want %d", tc.path, res.StatusCode, tc.status)
		}
	}
}
