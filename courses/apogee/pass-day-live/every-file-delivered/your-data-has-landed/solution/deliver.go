package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Customer is one customer's webhook receiver, as customers.json lists it.
type Customer struct {
	ID            string `json:"id"`             // "kessick"
	Name          string `json:"name"`           // "Kessick Polytechnic": its satellites' operator
	WebhookURL    string `json:"webhook_url"`    // where its events go
	WebhookSecret string `json:"webhook_secret"` // what they are signed with; never logged
}

// loadCustomers reads customers.json in dir strictly.
func loadCustomers(dir string) ([]Customer, error) {
	f, err := os.Open(filepath.Join(dir, "customers.json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cs []Customer
	if err := decodeStrict(f, &cs, "customers"); err != nil {
		return nil, fmt.Errorf("customers.json: %w", err)
	}
	return cs, nil
}

// Deliveries tells customers that their data has landed: one data.landed event per landed
// file, by bookingd's signed webhook.
type Deliveries struct {
	r      *run
	byName map[string]Customer // by name: a satellite's operator
	ctx    context.Context     // ends when the rehearsal stops
	wg     sync.WaitGroup
	mu     sync.Mutex
	made   int // events made so far
}

// newDeliveries returns the deliveries to the given customers; ctx ends them.
func newDeliveries(ctx context.Context, r *run, customers []Customer) *Deliveries {
	d := &Deliveries{r: r, byName: map[string]Customer{}, ctx: ctx}
	for _, c := range customers {
		d.byName[c.Name] = c
	}
	return d
}

// landed makes the event for a landed file and starts delivering it.
func (d *Deliveries) landed(l Landing, b Booking) {
	d.mu.Lock()
	d.made++
	id := fmt.Sprintf("evt-%06d", d.made)
	d.mu.Unlock()
	ev := Event{ID: id, Type: "data.landed", BookingID: l.BookingID, SatID: b.SatID, Station: b.Station,
		File: l.File, Bytes: l.Bytes, LandedAt: l.At}
	cust, ok := d.byName[satellite(d.r.n, b.SatID).Operator]
	d.wg.Go(func() {
		if !ok {
			sayAt(d.r.w, "UNDELIVERED %s", ev.File) // no customer operates the satellite
			return
		}
		d.deliver(cust, ev)
	})
}

// deliver sends one event to its customer's receiver: delivered when it answers 2xx.
func (d *Deliveries) deliver(cust Customer, ev Event) {
	status, err := d.post(cust, ev)
	if d.ctx.Err() != nil {
		return // stopping: not a failure
	}
	if err == nil && status/100 == 2 {
		sayAt(d.r.w, "DELIVERED %s", ev.File)
		return
	}
	sayAt(d.r.w, "UNDELIVERED %s", ev.File)
}

// post makes one attempt: the event as json.Marshal writes it, signed with this attempt's
// own timestamp. It returns the receiver's status, or an error when no answer came.
func (d *Deliveries) post(cust Customer, ev Event) (int, error) {
	body, err := json.Marshal(ev)
	if err != nil {
		return 0, err
	}
	ts := strconv.FormatInt(d.r.w.Now().Unix(), 10)
	mac := hmac.New(sha256.New, []byte(cust.WebhookSecret))
	mac.Write([]byte(ts + "."))
	mac.Write(body)
	req, err := http.NewRequestWithContext(d.ctx, http.MethodPost, cust.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Apogee-Timestamp", ts)
	req.Header.Set("Apogee-Signature", hex.EncodeToString(mac.Sum(nil)))
	res, err := d.r.w.HTTP.Do(req)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	return res.StatusCode, nil
}

// wait returns once every delivery has finished.
func (d *Deliveries) wait() { d.wg.Wait() }
