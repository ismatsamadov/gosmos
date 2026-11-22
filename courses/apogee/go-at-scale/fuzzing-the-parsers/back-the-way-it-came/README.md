---
slug: back-the-way-it-came
title: Back the Way It Came
type: code_tests
position: 6
difficulty: 9
topics: formatting, testing
created: 2026-11-11T12:51:43Z
modified: 2026-11-11T12:51:43Z
---
Joss's next flight software writes a battery of 7.91 V as `batt_v=+7.91` one time and `batt_v=07.91` another. tlmd accepts both and writes both back as `7.91`: two lines for one frame, and the archive can't tell them apart.

A round trip is a property every accepted line must have: format the frame `parse` returns, and you get the line back, byte for byte. A parser that accepts a second way of writing a frame fails it on that second way:

```go
v, _ := strconv.ParseFloat("+7.91", 64)
fmt.Printf("batt_v=%.2f\n", v) // batt_v=7.91
```

## Your task

`checkLine` checks the frame's values only. When `parse` accepts `line`, also fail `t` unless `formatLine` of the frame is exactly `line`.

> [!NOTE]
> Each of Mirela's four mutants lets one of Joss's lines through: `+7.91`, `07.91`, a time ending `+00:00`, and the sequence number `02214`.
