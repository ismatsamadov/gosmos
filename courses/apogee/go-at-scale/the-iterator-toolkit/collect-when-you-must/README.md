---
slug: collect-when-you-must
title: Collect When You Must
type: code
position: 0
difficulty: 8
topics: iterators, maps, slices, sorting
created: 2026-11-11T22:52:57Z
modified: 2026-11-12T08:43:35Z
---
The Board's station list still shows only the five stations it was written for, none of the first wave. Ingrid wants every live station on it, in code order.

Much of the standard library makes and takes sequences. `maps.Keys` yields a map's keys and `maps.Values` its values, in no fixed order. A sequence hands over one value at a time; when you need them all at once, as a sort does, gather them: `slices.Collect` makes a slice of any sequence, and `slices.Sorted` makes one and sorts it.

```go
gustKmh := map[string]int{"CWR": 81, "OST": 64}
kmh := slices.Sorted(maps.Values(gustKmh))
fmt.Println(kmh) // [64 81]
```

## Your task

`main` fills `live` with every station flying from Monday, keyed by its code, then prints the old list. Print the codes in `live` instead, sorted and joined by spaces:

```text
BRM CAL CWR LAN OST PEL SQT TIR VSK YST
```

> [!TIP]
> `maps` and `slices` are packages of their own: add both to the imports.
