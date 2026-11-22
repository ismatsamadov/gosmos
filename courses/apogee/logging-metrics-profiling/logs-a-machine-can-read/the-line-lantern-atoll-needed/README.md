---
slug: the-line-lantern-atoll-needed
title: The Line Lantern Atoll Needed
type: code_tests
position: 9
difficulty: 7
topics: errors, json, logging, time
created: 2026-11-07T07:51:37Z
modified: 2026-11-13T16:25:47Z
---
Friday 10 November. At 14:00 the new uplinkd goes to Lantern Atoll, and first Nkem replays the 02:02:39 failure on the lab bench: "Everything we didn't know at 02:48. In one line."

That line fits in two functions: a JSON logger whose records carry the clock you give it, and one record per pass with typed values, the wait in milliseconds and the error when there is one. `stamp`, in `clock.go`, wraps a handler so each record gets the time `now` returns:

```go
at := func() time.Time { return prep }
h := stamp(slog.NewTextHandler(w, nil), at)
slog.New(h).Info("pass started")
// time=2028-11-08T01:53:09.000Z level=INFO
//   msg="pass started"
```

## Your task

1. Make `newPassLog` return a logger on a JSON handler for `w`, wrapped by `stamp` with `now`.
2. Make `passLine` write one record: ERROR `pass failed` with `pass_id`, `max_elev_deg`, `msAttr("wait", wait)` and `err`; or, when `err` is nil, INFO `pass flown` with the same attributes but no `err`.
