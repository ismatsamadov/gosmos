---
slug: unknown-station-names
title: Unknown Station Names
type: code_tests
position: 6
difficulty: 4
topics: maps, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Hollis sends Furrowlight's bookings per station every week, as `{"YST":6,"SQT":4}`, and this week Cape Wren joins as `"CWR"`. No new build of gctl will reach this road, so the van must read whatever station codes arrive.

A JSON object decodes into a `map[string]int` whatever its keys. A value that is not a whole number, such as `2.5` or `"four"`, is an error. A key whose value is 0 is still in the map, and the two-value lookup tells it from a missing one.

```go
data := []byte(`{"stow_kmh":72,"gust_kmh":0}`)
var limits map[string]int
err := json.Unmarshal(data, &limits)
kmh, ok := limits["gust_kmh"]
fmt.Println(kmh, ok, err) // 0 true <nil>
```

## Your task

Implement `decodeCounts`: decode the object into a map and return it, keeping stations listed with 0. On an error, return nil with it, never the counts that did decode.
