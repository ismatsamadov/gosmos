---
slug: retry-without-hammering
title: Retry Without Hammering
type: code_tests
position: 2
difficulty: 10
topics: errors, http, time
created: 2026-11-12T05:36:26Z
modified: 2026-11-12T11:16:23Z
---
Friday 13 July. Yusra: "TL;DR: Cirrine's receiver restarts at 09:32 on Tuesday 24 July, eighteen seconds before Cirrine-5's file lands."

A receiver that fails may answer a moment later, so a failed delivery goes again after a wait: after failed attempt n, 2 s doubled n−1 times, at most 5 minutes, times a random factor from [0.5, 1.0), so senders don't all come back at once:

```text
after attempt 1: 1 s to under 2 s
after attempt 2: 2 s to under 4 s
after attempt 7: 64 s to under 128 s
```

## Your task

Retry on a transport error, 408, 429 and 5xx, at most 8 attempts, the same event ID, freshly timestamped and signed. Any other answer, or the 8th failure, makes the event undelivered: print `UNDELIVERED <file>`.

After the week, save `reports/undelivered.txt`, even empty: a line per undelivered event, by file, with its customer's ID and last answer, a status code or `error`:

```text
cirrine-8-20290726-0132.dat cirrine 503
```

The week ends once every delivery has; a stop ends every wait and is not a failure. The tests check Cirrine-5's file, receivers answering 503 or 400 or down, and a short pass at Lantern Atoll.
