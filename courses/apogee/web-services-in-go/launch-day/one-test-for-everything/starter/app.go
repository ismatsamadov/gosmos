package main

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// App is the assembled service: the booking service, its notifier, and the handler that serves them.
type App struct {
	Service  *Service
	Notifier *Notifier
	Handler  http.Handler // the mux inside every middleware
}

// newApp builds the whole service from cfg, with its bookings and its outbox in cfg.Dir.
func newApp(cfg Config) (*App, error) {
	svc, err := newService(cfg)
	if err != nil {
		return nil, err
	}
	outbox, err := OpenOutbox(filepath.Join(cfg.Dir, "outbox.json"))
	if err != nil {
		return nil, err
	}
	sender := &Sender{
		Client: &http.Client{
			Timeout:       5 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
		Customers: svc.Customers,
		Now:       cfg.Now,
		Rand:      cfg.Rand,
		Wait:      sleep,
	}
	logs := cfg.Log
	if logs == nil {
		logs = io.Discard
	}
	return &App{
		Service:  svc,
		Notifier: &Notifier{Store: svc.Store, Customers: svc.Customers, Sender: sender, Outbox: outbox},
		Handler:  svc.handler(log.New(logs, "", 0)),
	}, nil
}
