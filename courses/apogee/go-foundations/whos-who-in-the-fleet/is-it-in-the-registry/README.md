---
slug: is-it-in-the-registry
title: Is It in the Registry?
type: code_tests
position: 1
difficulty: 4
topics: maps
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Ask the registry for `brine-2` and you get an empty operator, which the dashboard took for a real satellite with a blank field. A missing ID has to say so.

Looking up a missing key gives the zero value, which can pass for real data. Ask for a second result and the map says whether the key is there. This is the comma-ok form:

```go
masks := map[string]float64{"SQT": 3}
mask, ok := masks["XYZ"]
fmt.Println(mask, ok) // 0 false
```

## Your task

Change `operatorOf` so it returns `"unknown satellite"` when `id` isn't in `reg`, and the stored operator when it is, even an empty one.

> [!NOTE]
> Pipit-3 went into the registry on launch day, before anyone typed its operator, so its operator is `""`.
