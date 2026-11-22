package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
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
// file, by bookingd's signed webhook and on its retry schedule.
type Deliveries struct {
	r      *run
	byName map[string]Customer // by name: a satellite's operator
	ctx    context.Context     // ends when the rehearsal stops
	wg     sync.WaitGroup
	mu     sync.Mutex
	made   int           // events made so far
	failed []undelivered // events that never will be delivered
}

// undelivered is an event that never reached its customer.
type undelivered struct {
	file, customer, last string // last: the last status code, or "error"
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
			d.give(ev, Customer{ID: "none"}, 0, "error", fmt.Errorf("no customer operates %s", b.SatID))
			return
		}
		d.deliver(cust, ev)
	})
}

// deliver sends one event until its customer's receiver takes it: retried on a transport
// error, 408, 429 and 5xx, at most 8 attempts, with growing waits between them.
func (d *Deliveries) deliver(cust Customer, ev Event) {
	for attempt := 1; ; attempt++ {
		status, err := d.post(cust, ev)
		if d.ctx.Err() != nil {
			return // stopping: not a failure
		}
		if err == nil && status/100 == 2 {
			sayAt(d.r.w, "DELIVERED %s", ev.File)
			d.r.met.Deliveries.With(cust.ID, "delivered").Inc()
			d.r.logs.delivered(ev.ID, cust.ID, attempt)
			return
		}
		retry := err != nil || status == 408 || status == 429 || status >= 500
		last := "error"
		if err == nil {
			last, err = strconv.Itoa(status), fmt.Errorf("receiver answered %d", status)
		}
		if !retry || attempt == 8 {
			d.give(ev, cust, attempt, last, err)
			return
		}
		if !sleep(d.ctx, backoff(attempt)) {
			return
		}
	}
}

// give records an event that never will be delivered, after its attempts, the last answer
// and its error.
func (d *Deliveries) give(ev Event, cust Customer, attempts int, last string, err error) {
	sayAt(d.r.w, "UNDELIVERED %s", ev.File)
	d.r.met.Deliveries.With(cust.ID, "undelivered").Inc()
	d.r.logs.undelivered(ev.ID, cust.ID, attempts, err)
	d.mu.Lock()
	d.failed = append(d.failed, undelivered{ev.File, cust.ID, last})
	d.mu.Unlock()
}

// backoff is the wait after failed attempt n: 2 s doubled for each earlier failure, at most
// 5 minutes, times a factor drawn from [0.5, 1.0).
func backoff(n int) time.Duration {
	d := 5 * time.Minute
	if n < 9 {
		d = min(d, 2*time.Second<<(n-1))
	}
	return time.Duration(float64(d) * (0.5 + rand.Float64()/2))
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

// writeUndelivered saves reports/undelivered.txt in dir: one line per event never
// delivered, by file: "<file> <customer> <last answer>".
func (d *Deliveries) writeUndelivered(dir string) error {
	d.mu.Lock()
	failed := slices.Clone(d.failed)
	d.mu.Unlock()
	slices.SortFunc(failed, func(a, b undelivered) int { return strings.Compare(a.file, b.file) })
	var buf bytes.Buffer
	for _, u := range failed {
		fmt.Fprintf(&buf, "%s %s %s\n", u.file, u.customer, u.last)
	}
	return saveWhole(filepath.Join(dir, "reports", "undelivered.txt"), buf.Bytes())
}
