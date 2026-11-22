---
slug: a-line-for-every-pass
title: A Line for Every Pass
type: code_tests
position: 0
difficulty: 9
topics: json, logging, time
created: 2026-11-12T09:13:51Z
modified: 2026-11-12T11:24:16Z
---
Sunniva's forecast: "Thursday 26th, storm. Gusts past 72 by 11:00. Fiber at risk." Ingrid wants that Thursday in the rehearsal's logs, as the night shift reads them.

`slog.NewJSONHandler` writes each record as one line of JSON; `slog.Group` nests attributes:

```go
l := slog.New(slog.NewJSONHandler(os.Stderr, nil))
l.Warn("dish stowed", "dish", 1,
	slog.Group("wind", "kmh", 81, "limit", 72))
```

## Your task

Write `logs/uplinkd.log` and `logs/dlinkd.log`, saved whole when the week ends, stopped or not: a record per event, its time from `w.Now()`, `program` the log's name, `site` the station or `HQ`, and these attributes, uplinkd's then dlinkd's:

```text
pass started  INFO   pass pass_id dish
pass flown    INFO   pass pass_id max_elev_deg
                     frames commands
pass failed   ERROR  pass pass_id max_elev_deg
                     reason err
pass skipped  WARN   pass pass_id reason
line lost     WARN   pass_id err
```

```text
file landed        INFO   file bytes
                          booking_id pass_id
event delivered    INFO   event_id customer
                          attempts
event undelivered  ERROR  event_id customer
                          attempts err
```

`pass` groups `sat_id`, `station`, `aos`, `los`; `frames` counts TLM frames with a good CRC, `commands` the pass's commands that ran, `customer` an ID, and a file's `pass_id` the pass it lands on. A stopped pass has no end record. The tests read every record and look for secrets.
