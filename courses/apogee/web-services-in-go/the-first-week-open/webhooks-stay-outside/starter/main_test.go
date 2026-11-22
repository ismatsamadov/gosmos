package main

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestCheckDial(t *testing.T) {
	cases := []struct {
		address string
		inside  bool
	}{
		{"203.0.113.7:443", false}, // stands for a receiver out on the internet
		{"198.51.100.20:8443", false},
		{"172.32.0.1:443", false}, // just past 172.31.255.255
		{"[2001:db8::7]:443", false},
		{"10.20.0.31:80", true}, // the weather service
		{"172.16.5.4:443", true},
		{"172.31.255.254:80", true},
		{"192.168.1.20:8080", true},
		{"127.0.0.1:8080", true},
		{"127.8.9.10:80", true},
		{"[::1]:443", true},
		{"169.254.169.254:80", true},
		{"[fe80::1]:80", true},
		{"[fd12:3456::7]:443", true},
		{"0.0.0.0:80", true},
		{"[::]:80", true},
		{"weather-service:80", true}, // not an address: refused, never guessed at
	}
	for _, c := range cases {
		err := checkDial("tcp", c.address, nil)
		switch {
		case c.inside && !errors.Is(err, errInside):
			t.Errorf("checkDial(\"tcp\", %q) = %v, want an error wrapping errInside", c.address, err)
		case !c.inside && err != nil:
			t.Errorf("checkDial(\"tcp\", %q) = %v, want nil: it is out on the internet", c.address, err)
		}
	}
}

func TestWebhookClientNeverDialsInside(t *testing.T) {
	var reached atomic.Int32
	weather := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached.Add(1)
		http.NotFound(w, r)
	}))
	defer weather.Close()
	_, port, err := net.SplitHostPort(weather.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	addresses := []struct{ name, url string }{
		{"the weather service's address, 127.0.0.1", weather.URL + "/apogee-events"},
		{"a name for it, localhost, resolved from this Mac's own hosts file", "http://localhost:" + port + "/apogee-events"},
	}
	for _, a := range addresses {
		resp, err := webhookClient().Post(a.url, "application/json", strings.NewReader(`{"id":"evt-000121","type":"data.landed","booking_id":"bk-4913"}`))
		if err == nil {
			resp.Body.Close()
		}
		if !errors.Is(err, errInside) {
			t.Errorf("webhookClient().Post to %s = %v, want an error wrapping errInside", a.name, err)
		}
	}
	if n := reached.Load(); n != 0 {
		t.Errorf("the weather service got %d requests, want 0: the check runs before any connection", n)
	}
}
