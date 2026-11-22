package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestRetryableAnswers(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{"Cirrine down for maintenance", http.StatusServiceUnavailable, true},
		{"a receiver that crashed", http.StatusInternalServerError, true},
		{"a proxy that lost its receiver", http.StatusBadGateway, true},
		{"a receiver too slow to read the request", http.StatusRequestTimeout, true},
		{"a receiver asking the service to slow down", http.StatusTooManyRequests, true},
		{"a request the receiver can't read", http.StatusBadRequest, false},
		{"a receiver that refuses the signature", http.StatusUnauthorized, false},
		{"an address nobody serves", http.StatusNotFound, false},
		{"Brinewatch's retired address", http.StatusGone, false},
		{"a receiver that rejects the event", http.StatusUnprocessableEntity, false},
		{"a receiver that moved for good", http.StatusMovedPermanently, false},
		{"a receiver that moved for now", http.StatusFound, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &AnswerError{StatusCode: tt.status, Status: strconv.Itoa(tt.status) + " " + http.StatusText(tt.status)}
			if got := retryable(tt.status, err); got != tt.want {
				t.Errorf("retryable(%d, %v) = %v, want %v", tt.status, err, got, tt.want)
			}
		})
	}
	if retryable(http.StatusNoContent, nil) {
		t.Error("retryable(204, nil) = true, want false: the event arrived")
	}
}

// noAnswer returns the error a client gets when a request to a stalled receiver ends without an answer.
func noAnswer(t *testing.T, ctx context.Context, timeout time.Duration) error {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := (&http.Client{Timeout: timeout}).Do(req)
	if err == nil {
		resp.Body.Close()
		t.Fatal("the stalled receiver answered, want no answer")
	}
	return err
}

func TestRetryableWithoutAnswer(t *testing.T) {
	gone := httptest.NewServer(http.NotFoundHandler())
	gone.Close()
	_, refused := (&http.Client{Timeout: 2 * time.Second}).Get(gone.URL)
	if refused == nil {
		t.Fatal("a closed receiver answered, want a refused connection")
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"a refused connection", refused, true},
		{"a receiver that never answered", noAnswer(t, context.Background(), 50*time.Millisecond), true},
		{"the service stopping", noAnswer(t, canceled, 2*time.Second), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := retryable(0, tt.err); got != tt.want {
				t.Errorf("retryable(0, %v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
