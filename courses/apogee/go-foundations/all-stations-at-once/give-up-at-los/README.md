---
slug: give-up-at-los
title: Give Up at LOS
type: code_tests
position: 7
difficulty: 7
topics: channels, select, time
created: 2026-11-01T07:36:47Z
modified: 2026-11-01T07:36:47Z
---
The uplink waits for Pipit-2 to acknowledge each command. If no acknowledgment comes before LOS, it must stop waiting: the satellite has set, and nothing more will arrive.

`time.After(d)` returns a channel that delivers one value once `d` has passed. As a case in a `select`, it puts a limit on how long the other cases may take.

```go
select {
case kmh := <-anemometer:
	fmt.Println("Cape Wren wind", kmh, "km/h")
case <-time.After(5 * time.Second):
	fmt.Println("anemometer silent: stowing")
}
```

## Your task

`waitAck` waits forever. Make it `select` between `acks` and a timer: return the acknowledgment and `nil` if one arrives, or `""` and `ErrNoAck` once `limit` has passed.
