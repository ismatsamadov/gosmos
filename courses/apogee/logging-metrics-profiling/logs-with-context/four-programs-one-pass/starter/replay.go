package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// replayNovember8 plays the night of 7 to 8 November through the four programs, as Nkem does on the
// lab bench. Each program's logger comes from programLog, on its own writer: bookingd, passplan, and
// uplinkd at headquarters (hq) and at Lantern Atoll (lan).
func replayNovember8(bookingd, passplan, hq, lan io.Writer) {
	// The bench's clock: at sets it to each record's time before the record is written.
	var now time.Time
	clock := func() time.Time { return now }
	at := func(t string) {
		var err error
		if now, err = time.Parse(time.RFC3339, t); err != nil {
			panic(err) // the replay's times are written out below: a bad one is a typo in the bench
		}
	}
	bk := programLog(bookingd, "bookingd", "HQ", clock)
	pp := programLog(passplan, "passplan", "HQ", clock)
	up := programLog(hq, "uplinkd", "HQ", clock)
	st := programLog(lan, "uplinkd", "LAN", clock)

	// The night's four passes at Lantern Atoll, by AOS, and the context of each one's work.
	utc := func(h, m, s int) time.Time { return time.Date(2028, 11, 8, h, m, s, 0, time.UTC) }
	brine3 := PassView{ID: "brine-3/LAN/14383", Pass: Pass{SatID: "brine-3", Station: "LAN", AOS: utc(1, 48, 20), LOS: utc(1, 54, 41), MaxElevDeg: 19.1}}
	cirrine1 := PassView{ID: "cirrine-1/LAN/16402", Pass: Pass{SatID: "cirrine-1", Station: "LAN", AOS: utc(1, 55, 9), LOS: utc(2, 2, 39), MaxElevDeg: 22.6}}
	furrowA := PassView{ID: "furrow-a/LAN/20998", Pass: Pass{SatID: "furrow-a", Station: "LAN", AOS: utc(2, 39, 45), LOS: utc(2, 46, 24), MaxElevDeg: 21.7}}
	brine3Later := PassView{ID: "brine-3/LAN/14384", Pass: Pass{SatID: "brine-3", Station: "LAN", AOS: utc(3, 26, 17), LOS: utc(3, 32, 10), MaxElevDeg: 16.6}}
	on := func(p PassView) context.Context { return withPass(context.Background(), p.ID) }
	// failed writes the line a pass ends with at LOS when its time runs out, wait after its preparation began.
	failed := func(p PassView, wait time.Duration) {
		err := fmt.Errorf("pass %s: %w", p.ID, context.DeadlineExceeded)
		st.LogAttrs(on(p), slog.LevelError, "pass failed",
			passGroup(p), slog.Float64("max_elev_deg", p.MaxElevDeg), msAttr("wait", wait), slog.Any("err", err))
	}

	// Monday 6 November: the customers' programs book the night's passes, each request with its ID in
	// its context, as bookingd's requestID middleware puts it there.
	for _, b := range []struct {
		at, id, request string
		p               PassView
	}{
		{"2028-11-06T09:41:27Z", "bk-5608", "req-004187", brine3},
		{"2028-11-06T09:41:28Z", "bk-5609", "req-004188", brine3Later},
		{"2028-11-06T10:14:03Z", "bk-5611", "req-004211", cirrine1},
		{"2028-11-06T15:20:55Z", "bk-5613", "req-004302", furrowA},
	} {
		at(b.at)
		ctx := context.WithValue(on(b.p), requestIDKey{}, b.request)
		bk.InfoContext(ctx, "booking created", "booking_id", b.id)
	}

	// Tuesday 7 November, just after midnight: passplan plans them, by AOS.
	at("2028-11-07T00:00:41Z")
	for _, b := range []struct {
		id string
		p  PassView
	}{{"bk-5608", brine3}, {"bk-5611", cirrine1}, {"bk-5613", furrowA}, {"bk-5609", brine3Later}} {
		pp.InfoContext(on(b.p), "pass planned", "booking_id", b.id)
	}

	// Wednesday 8 November: headquarters hands each pass's lines to the station, and the station flies.
	at("2028-11-08T01:31:05Z")
	up.InfoContext(on(brine3), "lines delivered", "status", 202)
	st.InfoContext(on(brine3), "lines queued", "lines", 2)
	at("2028-11-08T01:44:40Z")
	up.InfoContext(on(cirrine1), "command sealed", "seq", 4688)
	at("2028-11-08T01:46:20Z")
	st.InfoContext(on(brine3), "pass started", passGroup(brine3))
	at("2028-11-08T01:48:30Z")
	up.InfoContext(on(cirrine1), "lines delivered", "status", 202)
	st.InfoContext(on(cirrine1), "lines queued", "lines", 1)
	// A downpour: the radio loses Brine-3's carrier ten times a second, and the bench's clock counts whole seconds.
	at("2028-11-08T01:51:02Z")
	for range 10 {
		st.WarnContext(on(brine3), "carrier lost", "margin_db", -1.5)
	}
	at("2028-11-08T01:52:47Z")
	for range 5 {
		st.WarnContext(on(brine3), "carrier lost", "margin_db", -1.5)
	}
	at("2028-11-08T01:53:09Z")
	st.InfoContext(on(cirrine1), "pass started", passGroup(cirrine1))
	at("2028-11-08T01:54:41Z")
	failed(brine3, 8*time.Minute+21*time.Second)
	at("2028-11-08T02:02:39Z")
	failed(cirrine1, 9*time.Minute+30*time.Second)
	at("2028-11-08T02:37:45Z")
	st.InfoContext(on(furrowA), "pass started", passGroup(furrowA))
	at("2028-11-08T02:46:24Z")
	failed(furrowA, 8*time.Minute+39*time.Second)

	// 03:14: uplinkd at Lantern Atoll restarts, and the next pass flies.
	at("2028-11-08T03:18:40Z")
	up.InfoContext(on(brine3Later), "lines delivered", "status", 202)
	st.InfoContext(on(brine3Later), "lines queued", "lines", 2)
	at("2028-11-08T03:24:17Z")
	st.InfoContext(on(brine3Later), "pass started", passGroup(brine3Later))
	at("2028-11-08T03:27:38Z")
	st.LogAttrs(on(brine3Later), slog.LevelInfo, "pass flown", passGroup(brine3Later),
		slog.Float64("max_elev_deg", brine3Later.MaxElevDeg), msAttr("wait", 3*time.Minute+21*time.Second), msAttr("ack", 71*time.Second))
}
