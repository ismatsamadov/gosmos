---
slug: say-what-went-wrong
title: Say What Went Wrong
type: code_tests
position: 2
difficulty: 5
topics: errors, formatting
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
The log now says `battery too low` forty times a night. Ingrid wants to know which satellite, and how low.

`fmt.Errorf` builds an error the way `fmt.Sprintf` builds a string, with the same verbs. Write the message in lowercase with no final period, because it's often printed after other text:

```go
err := fmt.Errorf("%s: wind %.0f km/h",
	"CWR", 81.0)
fmt.Println("pass failed:", err)
// pass failed: CWR: wind 81 km/h
```

## Your task

Make `checkBattery` return an error like `thistle: battery 6.20 V below 6.60 V`: the satellite, the reading and `minBattV`, both voltages with two decimals. Healthy readings still return `nil`.
