---
slug: keys-without-formatting
title: Keys Without Formatting
type: code_tests
position: 6
difficulty: 8
topics: maps, performance, structs
created: 2026-11-12T03:10:33Z
modified: 2026-11-12T03:10:33Z
---
Thursday afternoon. Every hour the shadow run lists both planners' candidates for the day in one stream, to check that they see the same passes, and drops the repeats with a map keyed by `fmt.Sprintf`: four allocations a pass.

A struct whose fields are all comparable is comparable too, so it can be a map key, and making one allocates nothing: its fields are copied in place, with no string built.

```go
type gust struct {
	Station string
	Minute  int64
}
peak := map[gust]float64{}
k := gust{"CWR", at.Unix() / 60}
peak[k] = max(peak[k], kmh)
```

## Your task

`dedupe` keys its map with a formatted string. Key it by `passKey` instead: the satellite, the station and the AOS in Unix seconds. Every pass still comes out once, where it first came, and a lookup should allocate nothing.

> [!TIP]
> `time.Time` makes a poor key: `==` compares its location too, so two values for the same instant can be two keys.
