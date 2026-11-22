---
slug: four-programs-one-pass
title: Four Programs, One Pass
type: code_tests
position: 7
difficulty: 8
topics: context, json, logging, time
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Friday 24 November, on the lab bench: before the deploy at 14:00, Nkem replays Cirrine-1's pass of 8 November through bookingd, passplan, the program that plans each day's passes, and both uplinkds. One search must follow it through all four.

Each program's logger puts the week's pieces together: a JSON handler, `stamp` with the replay's clock, `passHandler` around both, and `With` for the program and its site. A logger for someone at a terminal is built the same way:

```go
var v slog.LevelVar // the debug port sets it
o := &slog.HandlerOptions{Level: &v}
h := passHandler{slog.NewTextHandler(w, o)}
term := slog.New(h).With("site", "LAN")
```

## Your task

Make `programLog` return a logger on a JSON handler for `w`, wrapped by `stamp` with `now`, then by `passHandler`, with `program` and then `site` on every record. The test runs the given `replayNovember8`, and `follow` and `merge` must find these lines:

```text
Mon 10:14:03 bookingd HQ  booking created
Tue 00:00:41 passplan HQ  pass planned
Wed 01:44:40 uplinkd  HQ  command sealed
Wed 01:48:30 uplinkd  HQ  lines delivered
Wed 01:48:30 uplinkd  LAN lines queued
Wed 01:53:09 uplinkd  LAN pass started
Wed 02:02:39 uplinkd  LAN pass failed
```
