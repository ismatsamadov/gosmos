---
slug: from-aos-to-los
title: From AOS to LOS
type: code_tests
position: 0
difficulty: 8
topics: context, goroutines, time
created: 2026-11-11T16:41:21Z
modified: 2026-11-12T10:46:08Z
---
On the Fishbowl's screen the rehearsal still prints its three lines and stops; none of the plan's 14,297 contacts flies yet.

Each contact flies on its own clock, from its AOS to its LOS, so contacts overlap: at different stations, and at one station's two dishes. The tests fly the week in simulated time, whose clock jumps to the next instant anything is waiting for.

## Your task

After the bookings, fly every contact. Print a line at its AOS and one at its LOS, each starting with `w.Now()` as `Mon 15:04:05`:

```text
Thu 16:03:46 AOS thistle-2/PEL/02136 dish 1
Thu 16:09:59 LOS thistle-2/PEL/02136 flown
```

Once the last LOS has passed, print the count and return nil:

```text
week done: 14297 flown, 0 failed, 0 skipped
```

The tests check every line's instant, and that `rehearse` returns at the last LOS, not before.

> [!WARNING]
> Measure every wait from `w.Now()`, as `c.AOS.Sub(w.Now())` does. In the tests the `time` package's clock starts in 2000, so waiting for a 2029 instant, with `time.Until` or a deadline, takes 29 years: the tests report the run still running.
