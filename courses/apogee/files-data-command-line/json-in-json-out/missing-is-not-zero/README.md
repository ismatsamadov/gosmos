---
slug: missing-is-not-zero
title: Missing Is Not Zero
type: code_tests
position: 7
difficulty: 4
topics: pointers, structs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
For one day of sweeps along the horizon, Paz will send gctl an update that takes the lowest elevation Salar Quieto works at, its mask, from 3° to 0°. Most updates say nothing about the mask, and gctl reads their silence as 0° too.

A missing key leaves a `float64` field at 0, the same as a key that says 0. A pointer field tells them apart: `json.Unmarshal` leaves it nil when the key is missing or `null`, and points it at the value otherwise, 0 included.

```go
var u struct {
    GustKmh *int `json:"gust_kmh"`
}
_ = json.Unmarshal([]byte(`{"gust_kmh":0}`), &u)
gust := u.GustKmh
fmt.Println(gust != nil, *gust) // true 0
```

## Your task

Make `MaskDeg` a `*float64`. Then have `decodeMaskUpdate` return `oldDeg` when the update has no mask, missing or `null`, and the update's mask otherwise, 0 included.
