---
slug: one-configuration-file
title: One Configuration File
type: code_tests
position: 0
difficulty: 7
topics: errors, structs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Monday morning, before the last 140 meters of the climb, Paz writes the settings for Thursday, calibration day. A misspelled field has to stop gctl down here, not at the top.

Start from a value that already holds the defaults and decode the file into it: a field the file names replaces its default, a field it leaves out keeps it, and an explicit `0` stays `0`. A field that has no default is checked after decoding:

```go
limits := Limits{MaxWindKmh: 72}
err := json.Unmarshal(data, &limits)
if err != nil {
	return Limits{}, err
}
```

## Your task

Implement `loadConfig`:

1. Start from the defaults: `Station` set to `SQT`, `MinElevDeg` to 3 and `Timeout` to 30 seconds.
2. Decode one object from `r` with a `json.Decoder` that refuses unknown fields, then require that a second `Decode` returns `io.EOF`.
3. Return an error unless `Day` parses with the layout `time.DateOnly`.
