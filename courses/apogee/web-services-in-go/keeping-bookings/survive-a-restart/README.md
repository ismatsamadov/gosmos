---
slug: survive-a-restart
title: Survive a Restart
type: code_tests
position: 8
difficulty: 8
topics: files, http
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Friday 17 December, the last deploy of the year. Yusra forwards Hollis's email with a TL;DR on top: "Furrowlight's January bookings must still be there on 3 January." Below it, in capitals: "WE HAVE PLANNED JANUARY AROUND THEM."

A constructor that opens files returns an error, so a service that can't read its bookings never starts. `newService` builds the service from a `Config`, whose `Dir` is the data folder; in tests, a `t.TempDir()`.

```go
func newPoll(cfg Config) (*Poll, error) {
	cache, err := openCache(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("poll: %w", err)
	}
	return &Poll{Cache: cache}, nil
}
```

## Your task

1. `newService` keeps the bookings in a `MemStore`, which forgets them at every restart. Open a `FileStore` at `bookings.json` in `cfg.Dir` instead, and return its error if it fails.
2. In `routes`, register `GET /bookings` with the read-only `listBookings`.
