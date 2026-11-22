---
slug: how-long-to-wait
title: How Long to Wait
type: code_tests
position: 5
difficulty: 6
topics: errors, testing, time
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Paz is raising the timeout in gctl's settings file from `30s` to `2m` for the long sweeps ahead. A mistyped value must stop gctl with an error, not turn into a wait of zero seconds.

`time.ParseDuration` reads text such as `30s`, `2m` or `1m30s` into a `time.Duration`, with units from `ns` to `h`. It also accepts `0s` and negative text such as `-5s`: those parse, but they make no sense as a wait, so check the value after parsing.

```go
d, err := time.ParseDuration("1h15m")
if err != nil {
	return err
}
fmt.Println(d.Minutes()) // 75
```

## Your task

Implement `positiveDelay`: parse `text` with `time.ParseDuration` and return the duration. Return 0 and an error when the text does not parse, and when the duration is zero or negative.
