---
slug: quick-check-simulated-time
title: Quick Check: Simulated Time
type: quiz
position: 3
difficulty: 8
topics: context, networking, time
created: 2026-11-11T16:41:21Z
modified: 2026-11-11T16:41:21Z
---
Thursday morning, before the stop button. Nkem asks how a week can pass in a second without lying about time.

This is the heart of one contact, with `w` the world, `c` the contact and `conn` its call:

```go
time.Sleep(c.AOS.Sub(w.Now()))
stamp := w.Now().Format("Mon 15:04:05")
untilLOS := c.LOS.Sub(w.Now())
pctx, cancel := context.WithTimeout(ctx, untilLOS)
defer cancel()
hangUp := func() { conn.Close() }
stopHangUp := context.AfterFunc(pctx, hangUp)
defer stopHangUp()
```

## Your task

Answer three questions: where a line's time comes from, why the pass's context gets a timeout, and what ends a read on a silent station.
