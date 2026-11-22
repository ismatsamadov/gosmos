package main

import (
	"bufio"
	"cmp"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

// logBook is the rehearsal's JSON logs, as the services write theirs: one record per event,
// stamped with the rehearsal's clock. Each is written into a temporary sibling as the run
// goes and renamed into place at its end, so it is saved whole.
type logBook struct {
	files    map[string]*logFile     // "uplinkd" ... by program
	stations map[string]*slog.Logger // uplinkd's, by station
	dlinkd   *slog.Logger
	pager    *slog.Logger
}

// logFile is one log on its way to disk. Every record is one Write.
type logFile struct {
	mu   sync.Mutex
	path string        // where it goes: logs/uplinkd.log
	tmp  *os.File      // its temporary sibling
	bw   *bufio.Writer // records not yet written to tmp
	err  error
}

// Write adds one record.
func (f *logFile) Write(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return 0, f.err
	}
	n, err := f.bw.Write(p)
	f.err = err
	return n, err
}

// newLogs starts the logs of a rehearsal at the given stations in dir, every record's time
// from w's clock.
func newLogs(w World, sites []Site) (*logBook, error) {
	b := &logBook{files: map[string]*logFile{}, stations: map[string]*slog.Logger{}}
	dir := filepath.Join(w.Dir, "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	opts := &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if len(groups) == 0 && a.Key == slog.TimeKey {
			return slog.Time(slog.TimeKey, w.Now())
		}
		return a
	}}
	for _, program := range []string{"uplinkd", "dlinkd", "pager"} {
		path := filepath.Join(dir, program+".log")
		tmp, err := os.CreateTemp(dir, "."+program+".log-*")
		if err != nil {
			b.discard()
			return nil, err
		}
		b.files[program] = &logFile{path: path, tmp: tmp, bw: bufio.NewWriterSize(tmp, 64<<10)}
	}
	logger := func(program, site string) *slog.Logger {
		return slog.New(slog.NewJSONHandler(b.files[program], opts)).With("program", program, "site", site)
	}
	for _, s := range sites {
		b.stations[s.Code] = logger("uplinkd", s.Code)
	}
	b.dlinkd = logger("dlinkd", "HQ")
	b.pager = logger("pager", "HQ")
	return b, nil
}

// discard removes the logs' temporary files.
func (b *logBook) discard() {
	for _, f := range b.files {
		f.tmp.Close()
		os.Remove(f.tmp.Name())
	}
}

// passGroup is the pass group of a contact's records.
func passGroup(c Contact) slog.Attr {
	return slog.Group("pass", "sat_id", c.SatID, "station", c.Station, "aos", c.AOS, "los", c.LOS)
}

// started writes a contact's "pass started".
func (b *logBook) started(c Contact) {
	b.stations[c.Station].Info("pass started", passGroup(c), "pass_id", c.ID, "dish", c.Dish)
}

// ended writes how a contact ended: "pass flown", "pass failed" or "pass skipped".
func (b *logBook) ended(c Contact, o Outcome, p *pass) {
	l := b.stations[c.Station]
	switch o.Kind {
	case "flown":
		l.Info("pass flown", passGroup(c), "pass_id", c.ID, "max_elev_deg", c.MaxElevDeg, "frames", p.frames, "commands", p.up.answered)
	case "failed":
		l.Error("pass failed", passGroup(c), "pass_id", c.ID, "max_elev_deg", c.MaxElevDeg, "reason", o.Reason, "err", errText(o.Err))
	case "skipped":
		l.Warn("pass skipped", passGroup(c), "pass_id", c.ID, "reason", o.Reason)
	}
}

// errText is an error's text, "" for none.
func errText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// lineLost writes a "line lost".
func (b *logBook) lineLost(c Contact, err error) {
	b.stations[c.Station].Warn("line lost", "pass_id", c.ID, "err", errText(err))
}

// landed writes a "file landed".
func (b *logBook) landed(file string, bytes int, bookingID, passID string) {
	b.dlinkd.Info("file landed", "file", file, "bytes", bytes, "booking_id", bookingID, "pass_id", passID)
}

// landFailed writes a file that landed but could not be saved.
func (b *logBook) landFailed(file string, err error) {
	b.dlinkd.Error("file not saved", "file", file, "err", errText(err))
}

// delivered writes an "event delivered".
func (b *logBook) delivered(eventID, customer string, attempts int) {
	b.dlinkd.Info("event delivered", "event_id", eventID, "customer", customer, "attempts", attempts)
}

// undelivered writes an "event undelivered".
func (b *logBook) undelivered(eventID, customer string, attempts int, err error) {
	b.dlinkd.Error("event undelivered", "event_id", eventID, "customer", customer, "attempts", attempts, "err", errText(err))
}

// alert writes an "alert sent".
func (b *logBook) alert(a Alert) {
	b.pager.Info("alert sent", "rule", a.Rule, "station", a.Station, "kind", a.Kind, "value", a.Value, "to", a.To)
}

// save puts logs/uplinkd.log, logs/dlinkd.log and logs/pager.log in place.
func (b *logBook) save() error {
	var first error
	for _, f := range b.files {
		f.mu.Lock()
		err := f.err
		if err == nil {
			err = f.bw.Flush()
		}
		if cerr := f.tmp.Close(); err == nil {
			err = cerr
		}
		if err == nil {
			err = os.Rename(f.tmp.Name(), f.path)
		}
		if err != nil {
			os.Remove(f.tmp.Name())
			first = cmp.Or(first, err)
		}
		f.mu.Unlock()
	}
	return first
}
