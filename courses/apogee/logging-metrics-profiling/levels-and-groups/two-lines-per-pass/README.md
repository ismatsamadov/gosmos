---
slug: two-lines-per-pass
title: Two Lines per Pass
type: code_tests
position: 8
difficulty: 7
topics: logging, switch
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Friday 17 November: the levels go to Lantern Atoll at 14:00, and Ingrid wants a night at INFO she can read at a glance. Nkem: "Two lines a pass. Where it starts, how it ends. Trouble is WARN or ERROR, and every step between waits at DEBUG."

An event's level is chosen once, in one place, so every call that writes it agrees. A `switch` on the event's message picks the level, and a `case` may list several values:

```go
switch satID {
case "pipit-1", "pipit-2", "pipit-3":
	return "Apogee"
case "furrow-a", "furrow-b":
	return "Furrowlight"
}
return "unknown"
```

## Your task

Make `levelFor` return each uplinkd event's level from this list; any event it doesn't know stays at INFO:

```text
DEBUG  worker taken, lines sent,
       frame received, ack received
WARN   carrier lost, pass skipped
ERROR  pass failed
INFO   everything else
```

> [!NOTE]
> The test replays the night of 8 November, and at INFO each of its four passes must read as two lines, beside the warnings.
