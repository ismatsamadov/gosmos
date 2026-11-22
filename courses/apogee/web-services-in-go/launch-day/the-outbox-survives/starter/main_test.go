package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// thursday returns two of Thursday's deliveries: Cirrine's, refused twice and due again at 14:11:03, and
// Brinewatch's, not tried yet.
func thursday(cirrineID, brinewatchID string) (Delivery, Delivery) {
	cirrine := Delivery{
		Customer: "cirrine",
		Event: Event{ID: cirrineID, Type: "data.landed", BookingID: "bk-4861", SatID: "cirrine-1", Station: "LAN",
			File: "cirrine-1-20280127-1403.tif", Bytes: 505000000, LandedAt: time.Date(2028, 1, 27, 14, 10, 58, 0, time.UTC)},
		Attempts: 2,
		NextAt:   time.Date(2028, 1, 27, 14, 11, 3, 0, time.UTC),
	}
	brinewatch := Delivery{
		Customer: "brinewatch",
		Event: Event{ID: brinewatchID, Type: "data.landed", BookingID: "bk-4858", SatID: "brine-3", Station: "LAN",
			File: "brine-3-20280127-1346.tif", Bytes: 58750000, LandedAt: time.Date(2028, 1, 27, 13, 54, 52, 0, time.UTC)},
		NextAt: time.Date(2028, 1, 27, 13, 54, 52, 0, time.UTC),
	}
	return cirrine, brinewatch
}

// open opens the outbox at path, failing the test on an error.
func open(t *testing.T, path string) *Outbox {
	t.Helper()
	o, err := OpenOutbox(path)
	if err != nil {
		t.Fatalf("OpenOutbox(%s) = %v, want nil", filepath.Base(path), err)
	}
	return o
}

// newID takes the next event ID, failing the test on an error.
func newID(t *testing.T, o *Outbox) string {
	t.Helper()
	id, err := o.NewEventID()
	if err != nil {
		t.Fatalf("NewEventID() = %v, want nil", err)
	}
	return id
}

// asJSON writes deliveries as JSON, so two lists compare field by field, times included.
func asJSON(deliveries []Delivery) string {
	data, _ := json.Marshal(deliveries)
	return string(data)
}

func TestOutboxKeepsWhatIsPending(t *testing.T) {
	path := filepath.Join(t.TempDir(), "outbox.json")
	o := open(t, path)
	cirrine, brinewatch := thursday(newID(t, o), newID(t, o))
	if err := o.Put(cirrine); err != nil {
		t.Fatalf("Put(Cirrine's delivery) = %v, want nil", err)
	}
	if err := o.Put(brinewatch); err != nil {
		t.Fatalf("Put(Brinewatch's delivery) = %v, want nil", err)
	}
	reopened := open(t, path)
	if got, want := asJSON(reopened.Pending()), asJSON([]Delivery{cirrine, brinewatch}); got != want {
		t.Errorf("Pending() after a restart = %s, want %s", got, want)
	}
	if id := newID(t, reopened); id != "evt-000003" {
		t.Errorf("NewEventID() after a restart = %q, want %q", id, "evt-000003")
	}
	data, _ := os.ReadFile(path)
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil || fields["next_event"] == nil || fields["pending"] == nil {
		t.Errorf("outbox.json = %s, want an object with next_event and pending", data)
	}
}

func TestDeliveredEventLeavesTheFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "outbox.json")
	o := open(t, path)
	cirrine, brinewatch := thursday(newID(t, o), newID(t, o))
	o.Put(cirrine)
	o.Put(brinewatch)
	if err := o.Remove(brinewatch.Event.ID); err != nil {
		t.Fatalf("Remove(%s) = %v, want nil", brinewatch.Event.ID, err)
	}
	if got, want := asJSON(open(t, path).Pending()), asJSON([]Delivery{cirrine}); got != want {
		t.Errorf("Pending() after Brinewatch's event was delivered and the service restarted = %s, want %s", got, want)
	}
}

func TestEventIDsNeverRepeat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "outbox.json")
	if err := os.WriteFile(path, []byte(`{"next_event":117,"pending":[]}`), 0600); err != nil {
		t.Fatal(err)
	}
	o := open(t, path)
	id := newID(t, o)
	if id != "evt-000117" {
		t.Errorf("NewEventID() from next_event 117 = %q, want %q", id, "evt-000117")
	}
	cirrine, _ := thursday(id, "")
	o.Put(cirrine)
	o.Remove(id)
	if next := newID(t, open(t, path)); next != "evt-000118" {
		t.Errorf("NewEventID() after a restart with nothing pending = %q, want %q: an ID is never used twice", next, "evt-000118")
	}
}

func TestMissingOutbox(t *testing.T) {
	path := filepath.Join(t.TempDir(), "outbox.json")
	o := open(t, path)
	if got := o.Pending(); len(got) != 0 {
		t.Errorf("Pending() of a new outbox = %v, want nothing", got)
	}
	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("OpenOutbox(a missing file) made the file (%v), want nothing written before the first change", err)
	}
	if id := newID(t, o); id != "evt-000001" {
		t.Errorf("NewEventID() of a new outbox = %q, want %q", id, "evt-000001")
	}
}

func TestBrokenOutbox(t *testing.T) {
	path := filepath.Join(t.TempDir(), "outbox.json")
	if err := os.WriteFile(path, []byte(`{"next_event":118,"pending":[{"customer":`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenOutbox(path); err == nil {
		t.Errorf("OpenOutbox(a file cut short) = nil error, want an error: a broken outbox stops the start")
	}
}
