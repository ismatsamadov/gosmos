---
slug: do-not-hide-the-bug
title: Do Not Hide the Bug
type: code_tests
position: 11
difficulty: 6
topics: defer, errors, interfaces
created: 2026-11-04T17:00:51Z
modified: 2026-11-06T16:51:28Z
---
Mirela reviews the decoder boundary: "You caught it. Did you know what it was?" The older decoder signals a known internal fault with `decoderFault`; an unrelated programming bug must still be visible to the outer service boundary.

Recover narrowly. After recovery, use a type assertion to recognize the panic you can handle. Pass every other panic on with `panic(problem)`, preserving its original value. Never turn an unknown bug into a successful zero reading.

```go
var problem any = "heater stuck"
s, ok := problem.(string) // "heater stuck", true
n, isInt := problem.(int) // 0, false
```

## Your task

Change the deferred function in `runKnownDecoder`. Convert only a `decoderFault` into `fmt.Errorf("decoder fault: %s", fault.reason)`. Re-panic every other value unchanged. Successful calls still return their reading and nil error.
