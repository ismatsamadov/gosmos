---
slug: nothing-allocated-per-frame
title: Nothing Allocated per Frame
type: code_tests
position: 2
difficulty: 9
topics: conversions, profiling, slices
created: 2026-11-10T15:36:39Z
modified: 2026-11-10T15:36:39Z
---
Tuesday. During a pass, Lantern Atoll's uplinkd stores up to 500 frames a second, the numbered pieces of data the radio hears, each under its pass ID. It builds the ID with `fmt.Sprintf` for every frame: 500 new strings a second for the garbage collector to clear.

`append` and `strconv.AppendInt` write into a slice you pass in, and while it has room they allocate nothing. `testing.AllocsPerRun(n, f)` runs `f` n times and returns how many allocations a run made, on average:

```go
b := make([]byte, 0, 16)
n := testing.AllocsPerRun(100, func() {
	b = strconv.AppendInt(b[:0], 4688, 10)
})
fmt.Println(string(b), n) // 4688 0
```

## Your task

`appendPassID` appends what `fmt.Sprintf` makes. Make it append the satellite, `/`, the station, `/` and the orbit number padded with zeros to five digits, using only `append` and `strconv.AppendInt`.

> [!NOTE]
> Padding only adds zeros in front: an orbit number of six digits keeps all six.
