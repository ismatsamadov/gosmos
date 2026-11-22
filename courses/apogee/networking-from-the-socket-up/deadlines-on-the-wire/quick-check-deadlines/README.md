---
slug: quick-check-deadlines
title: Quick Check: Deadlines
type: quiz
position: 2
difficulty: 7
topics: networking, time
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T07:44:28Z
---
Wednesday. Sunniva's questions arrive over the backhaul, weather first, as ever: "Wind 70, falling. Deadline: instant or span? How removed? Line usable after?" They're about this code, where `t0` is the moment it starts:

```go
t0 := time.Now()
conn.SetReadDeadline(t0.Add(5 * time.Second))
f1, err := readFrame(br) // returns at t0+2s
time.Sleep(4 * time.Second)
f2, err := readFrame(br) // starts at t0+6s
```

## Your task

Answer Sunniva's three questions: when the second read gives up, how to remove a deadline, and what the connection can do after a read times out.
