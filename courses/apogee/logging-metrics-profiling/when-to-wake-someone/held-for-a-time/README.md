---
slug: held-for-a-time
title: Held for a Time
type: code_tests
position: 1
difficulty: 8
topics: methods, metrics, time
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
Ingrid's new `disk-full` rule fires at 0.90, but only after 15 minutes there. Lantern Atoll's disk can touch 91% for four minutes while a pass's recording lands, until headquarters collects it, and that shouldn't wake anyone; fifteen minutes should.

Every minute the pager takes a sample, the rule's value at that moment. The value at which a rule's problem starts is its threshold, `FireAt`, and `For` is how long the value must stay at or above it. Until then the rule is pending, since the first sample that reached it; a sample below `FireAt` makes it inactive, and the wait starts over. `Observe` hands each sample to `beforeFiring` until the rule fires:

```go
t0 := time.Date(2029, 1, 3, 9, 0, 0, 0, time.UTC)
t := t0.Add(14*time.Minute + 59*time.Second)
fmt.Println(t.Sub(t0) >= 15*time.Minute)
// false: still pending
```

## Your task

`beforeFiring` fires at the first sample at `FireAt`. Keep the time the value reached `FireAt` in `st.since`, fire only once `s.At` is at least `r.For` after it, and set `st.since` back to zero when a sample falls below `FireAt`.

> [!NOTE]
> `failed-pass` will have a `For` of 0, and a rule like that fires at its first sample.
