---
slug: the-rule-that-counts
title: The Rule That Counts
type: code_tests
position: 4
difficulty: 4
topics: functions, variables
created: 2026-11-04T16:56:59Z
modified: 2026-11-05T08:05:38Z
---
Ingrid asks how often each customer's checklist runs. A shared counter would mix the customers together; each checklist needs a count of its own.

A function literal can use a variable from its surrounding function. The returned function keeps that variable alive after the surrounding call ends: this is a closure. Calls to the same closure share its variable; each new call of the surrounding function makes a separate one.

```go
message := "ready"
announce := func() { fmt.Println(message) }
message = "listening"
announce() // listening
```

## Your task

In `makeCounter`, create a count before the returned function. Inside that function, increase the count and return it. Each new counter starts at 1 on its first call, then returns 2, 3 and so on independently.
