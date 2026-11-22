---
slug: quick-check-who-closes
title: Quick Check: Who Closes
type: quiz
position: 2
difficulty: 7
topics: channels
created: 2026-11-10T15:39:51Z
modified: 2026-11-10T15:39:51Z
---
Before the drill's twenty feeds start, Nkem draws three stages and four arrows on the whiteboard: a feed into `read`, `read` into `parse`, `parse` into `newest`, and `newest` out to the Board.

```go
lines := read(ctx, "PEL", feed)
frames := parse(ctx, lines)
health := newest(frames)
```

`read` returns a `<-chan Heard`, `parse` returns a `<-chan TelemetryFrame`, and `newest` ranges over `frames`.

## Your task

Answer Nkem's three questions: who closes `frames`, what adding `close(lines)` to this code would do, and when the loop in `newest` ends.
