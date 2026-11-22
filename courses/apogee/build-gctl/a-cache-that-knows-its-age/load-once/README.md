---
slug: load-once
title: Load Once
type: code_tests
position: 0
difficulty: 6
topics: functions, structs
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T17:14:45Z
---
A message from Sunniva at Ystmark, whose link home, the backhaul, now runs through a slow satellite: "Storm. Fiber down. Nine seconds a command." gctl reads the night's schedule again for every command, when once is enough.

`cache.go` brings `passCache`, which holds a loader, `load`: a function value it can call when it needs the passes, and not before. Its clock and time to live wait for the next lesson.

## Your task

1. Write `passes`: call `c.load` the first time, keep what it returns, and return the kept passes from then on. When `load` fails, return its error and keep nothing, so the next call tries again.
2. In `Shell`, replace the `passes` field with a `cache` field, a `*passCache` made in `newShell` with `loadTonight` from `data.go` as its loader.
3. In `cmdPasses`, take the passes from `s.cache.passes()`, and return its error if it has one.
