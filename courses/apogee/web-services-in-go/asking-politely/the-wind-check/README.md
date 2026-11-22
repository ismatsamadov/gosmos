---
slug: the-wind-check
title: The Wind Check
type: code_tests
position: 8
difficulty: 6
topics: context, errors, http, time
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Friday evening. Before the night's passes the poll checks the wind at every station, and Ingrid wants each answer in time: "GO, NO-GO, or the reason why. Nothing else."

`windCall` puts the week together. The client gives up after its `Timeout`, `Reading` returns a `*StatusError` for any status but 200, and a reading is fresh for 15 minutes after its `at`, exactly 15 included. Two formats you need:

```go
fmt.Printf("%.0f km/h\n", 61.0) // 61 km/h
fmt.Println(47 * time.Minute)   // 47m0s
```

## Your task

Write `windCall`: call `Reading` with `ctx` and `code`, and return the first of these lines that applies, from the top:

```text
NO-GO: mast silent         a 503
NO-GO: no weather          any other failure
NO-GO: reading 47m0s old   over 15 minutes old
NO-GO: wind 74 km/h, stow  from 72 km/h
GO: wind 61 km/h, warning  from 50 km/h
GO: wind 22 km/h           under 50 km/h
```

The age is `now` minus the reading's `at`, and the wind is `WindKmh`.
