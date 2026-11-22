---
slug: removing-brine-2
title: Removing Brine-2
type: code
position: 4
difficulty: 3
topics: maps
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Yusra's message opens, as always, with its summary: "TL;DR: Brine-2 never left its deployer. Take it out of the registry today."

The built-in `delete(m, key)` removes a key and its value from a map. Deleting a key that isn't there does nothing, so there's no need to check first.

```go
masks := map[string]float64{"YST": 5, "XYZ": 0}
delete(masks, "XYZ")
fmt.Println(len(masks)) // 1
```

## Your task

`loadRegistry` in `registry.go` returns the old fleet list, Brine-2 included, and that file is read-only. In `main`, delete `brine-2` from `reg` before it's checked, so the program prints:

```text
satellites: 8
brine-2 listed: false
```
