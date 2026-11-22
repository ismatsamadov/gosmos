---
slug: which-payload
title: Which Payload Is This?
type: code
position: 5
difficulty: 5
topics: interfaces, switch, types
created: 2026-11-01T07:36:46Z
modified: 2026-11-04T17:07:34Z
---
The downlink hands everything it receives to one function as an `any`: Furrow-A's images, routine equipment-status readings called housekeeping frames, and chunks of files. Each kind goes somewhere different.

A type switch asks which concrete type an interface value holds. `item.(type)` works only in a `switch`, and the case that names the value's type runs:

```go
switch reading.(type) {
case int:
	fmt.Println("millivolts")
case float64:
	fmt.Println("volts")
default:
	fmt.Println("unreadable")
}
```

## Your task

Write a type switch in `route`: an `Image` goes to `"imagery"`, `Housekeeping` to `"telemetry"`, a `Chunk` to `"reassembly"`, and anything else to `"quarantine"`.
