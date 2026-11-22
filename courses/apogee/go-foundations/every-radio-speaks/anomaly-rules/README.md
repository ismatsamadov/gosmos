---
slug: anomaly-rules
title: Anomaly Rules
type: code_tests
position: 7
difficulty: 6
topics: interfaces, methods, slices
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:26Z
---
Mirela wants to add anomaly rules without opening the pager's code. The pager keeps its rules in one `[]Rule`, and each rule type brings its own `Check` method.

A slice of an interface type holds values of different types side by side, and a method call on each element runs that value's own method:

```go
type Poller interface{ Ready() bool }

polls := []Poller{Wind{Kmh: 64}, Heater{On: true}}
for _, p := range polls {
	// runs Wind's Ready, then Heater's
	fmt.Println(p.Ready())
}
```

## Your task

Write the two `Check` methods: a rule that fires returns its message and `true`, otherwise `"", false`. `BatteryLow` fires below `MinV` with a message like `battery 6.40 V below 6.60 V`, and `TooHot` fires above `MaxC` with one like `temperature 51.0 C above 45.0 C`: the reading, then the rule's limit. A reading exactly at a limit doesn't fire.
