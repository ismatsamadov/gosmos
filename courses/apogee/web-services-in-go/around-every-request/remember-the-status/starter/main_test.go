package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusRecorder(t *testing.T) {
	tests := []struct {
		name   string
		answer func(w http.ResponseWriter)
		status int // what Status should record, and the writer inside receive
		bytes  int
	}{
		{"a booking that doesn't exist", func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, `{"error":"not_found"}`)
		}, http.StatusNotFound, 21},
		{"a body with no status first", func(w http.ResponseWriter) {
			io.WriteString(w, `{"status":"ok"}`)
		}, http.StatusOK, 15},
		{"a body in two writes", func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusCreated)
			io.WriteString(w, `{"id":"bk-4512",`)
			io.WriteString(w, `"sat_id":"cirrine-1"}`)
		}, http.StatusCreated, 37},
		{"a second status, too late", func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusCreated)
			w.WriteHeader(http.StatusInternalServerError)
		}, http.StatusCreated, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inner := httptest.NewRecorder()
			rec := &statusRecorder{ResponseWriter: inner}
			tt.answer(rec)
			if rec.Status != tt.status {
				t.Errorf("statusRecorder.Status = %d, want %d", rec.Status, tt.status)
			}
			if rec.Bytes != tt.bytes {
				t.Errorf("statusRecorder.Bytes = %d, want %d", rec.Bytes, tt.bytes)
			}
			if inner.Code != tt.status || inner.Body.Len() != tt.bytes {
				t.Errorf("the writer inside got %d and %d bytes, want %d and %d bytes", inner.Code, inner.Body.Len(), tt.status, tt.bytes)
			}
		})
	}
}

func TestStatusRecorderUnwrap(t *testing.T) {
	inner := httptest.NewRecorder()
	rec := &statusRecorder{ResponseWriter: inner}
	if err := http.NewResponseController(rec).Flush(); err != nil {
		t.Fatalf("http.NewResponseController(statusRecorder).Flush() = %v, want nil", err)
	}
	if !inner.Flushed {
		t.Errorf("after Flush, the writer inside: Flushed = %v, want %v", inner.Flushed, true)
	}
}
