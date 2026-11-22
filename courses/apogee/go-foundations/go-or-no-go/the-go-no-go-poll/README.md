---
slug: the-go-no-go-poll
title: The Go/No-Go Poll
type: code_tests
position: 7
difficulty: 5
topics: conditionals, switch
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Furrow-A rises over Salar Quieto in twenty minutes for its evening pass, and tonight your poll makes the call. Ingrid wants one reason, the first in checklist order, so the Fishbowl fixes the most basic problem first.

When several checks can fail, their order decides which failure you report. Check in order and return at the first failure; code that gets past every check can say yes. A `switch` with no value works the same way, since its first true case wins.

```go
if !sent {
	return "not sent"
}
if !heard {
	return "no echo"
}
return "echo received"
```

## Your task

Make `poll` return the first of these that applies, in this order: `NO-GO: station down` if the station isn't up, `NO-GO: wind` from 72 km/h, `NO-GO: below mask` if the pass peaks under Salar Quieto's 3° mask, and `NO-GO: battery` under 6.6 V. Otherwise, return `GO`.
