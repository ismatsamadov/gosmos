---
slug: weather-once-a-minute
title: Weather Once a Minute
type: code_tests
position: 0
difficulty: 8
topics: maps, mutexes, time
created: 2026-11-07T06:23:24Z
modified: 2026-11-07T06:23:24Z
---
Thursday. Tomás asks why his old weather service gets twenty calls a minute from the Board alone: five stations, four refreshes a minute. Nkem routed the Board's weather through `cache.go` this morning, and its `reading` still asks every time.

A cache keeps each answer with the time it came and hands it back while it is young. Its clock is the Board's, `c.now`, so the tests move time without waiting. A reading is kept for `c.ttl`, a minute: at exactly 60 seconds it is fetched again. An error is never kept, so a silent mast is asked again at the next refresh, and each station's reading keeps its own time.

## Your task

Write `reading`:

1. Under `c.mu`, look up the station in `c.kept`; a reading kept less than `c.ttl` ago is returned as it is.
2. Otherwise ask `c.fetch`. Keep a reading it returns, with the time it came, and return it; return an error without keeping anything.

Four refreshes 15 seconds apart then ask once per station, and a fifth, a minute after the first, asks again.
