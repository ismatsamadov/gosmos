---
slug: say-it-once-then-count
title: Say It Once, Then Count
type: code_tests
position: 6
difficulty: 8
topics: logging, mutexes, time
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Through a downpour one line repeats ten times a second, `carrier lost`, the radio's warning that the satellite's signal is too weak to follow, and buries every other line of the night.

Sampling keeps a flood readable: of each level and message, the sampler passes the first 10 records of a UTC minute, then one in a hundred, the 110th, the 210th and so on, and every ERROR record. The next one it passes after drops carries `dropped`, how many went since the last one passed, from one minute into the next. A handler can add to a record before it hands it on:

```go
func (h tagged) Handle(ctx context.Context,
	r slog.Record) error {
	r.AddAttrs(slog.String("site", "LAN"))
	return h.inner.Handle(ctx, r)
}
```

## Your task

1. Make `Handle` pass ERROR records, and any level above, straight to `s.inner`.
2. Count each other record, under `s.shared.mu`, in the tally `s.shared.tallyOf(r)` returns. Keep it when its count in the minute is up to 10, or 110, 210 and so on; count each one you drop in the tally's `dropped` and in `s.shared.total`.
3. Add `dropped` to a record you keep after drops, then pass it on.
