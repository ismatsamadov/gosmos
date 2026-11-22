---
slug: quick-check-interfaces
title: Quick Check: Interfaces
type: quiz
position: 2
difficulty: 4
topics: interfaces
created: 2026-11-01T07:36:46Z
modified: 2026-11-22T10:08:26Z
---
While the Ka-band radio warms up on La Grande, Tomás leans over your desk with three quick questions. They're about these lines, where `XBandRadio` and `UHFRadio` are the radio types from your pass routine and everything builds:

```go
type Radio interface {
	RateMbps() float64
}

var r Radio = XBandRadio{}
r = UHFRadio{}
var v any = 437.25
```

## Your task

Pick one answer for each question, then press **Submit**.
