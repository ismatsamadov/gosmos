---
slug: rotation-day-hour-by-hour
title: Rotation Day, Hour by Hour
type: code_tests
position: 7
difficulty: 10
topics: crypto, sorting, time
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T19:06:42Z
---
14:11:52. Furrow-A sets over Salar Quieto, the last of Rotation Day's ten passes. Before Mirela signs the log, Ingrid wants the whole day run through your code, from the Cirrine-1 pass that crossed midnight, and its report to match hers.

`day.go` holds the day: each pass with the times its PINGs were sealed, uplinkd's ring, each station's old and new certificate, and list 1. `Connect` makes a station's handshake and `Ping` delivers a sealed PING. Each pass becomes one line:

```text
04:12:05 pipit-1 PEL flown k-2028-10 1102
09:14:20 brine-3 LAN missed: <first error>
```

## Your task

Write `runDay`. For each slot, in AOS order:
1. `Connect` at AOS with the station's certificate in force, `New` once AOS reaches its `switchTime` from `d.SwitchFrom`, else `Old`, and with `d.List` once `d.ListAt` has come.
2. For each `SealAt`, `Ping` 5 seconds later, under the key `signingKey` picks at that `SealAt`.
3. Count the pass flown, with its key IDs joined by `+` and the serial, or missed at its first error.
