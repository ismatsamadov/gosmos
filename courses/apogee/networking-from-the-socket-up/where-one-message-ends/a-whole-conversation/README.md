---
slug: a-whole-conversation
title: A Whole Conversation
type: code_tests
position: 8
difficulty: 7
topics: errors, networking, streams
created: 2026-11-07T07:44:29Z
modified: 2026-11-07T07:44:29Z
---
Friday 17 March: dishd's first full session with the simulator, over Sunniva's slow test line. Nkem's one question: "Does every answer reach the request it belongs to?"

A controller answers frames in the order they arrive, so a client can send several requests at once and read the answers back in that order, one per request. An error from `recv` ends the session, a `*ByeError` included: the controller hanging up is never an answer. Return the answers you have, with the error:

```go
var temps []float64
for _, s := range readings {
	t, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return temps, err
	}
	temps = append(temps, t)
}
```

## Your task

`exchange` returns `nil`. Send every frame with `writeFrame`, then call `recv` once per frame and return the answers in order. On the first error, return the answers read so far with that error.
