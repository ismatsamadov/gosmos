---
slug: a-stream-that-stays-open
title: A Stream That Stays Open
type: code_tests
position: 2
difficulty: 8
topics: context, http, streams
created: 2026-11-07T06:23:23Z
modified: 2026-11-07T06:23:23Z
---
Fifteen seconds is too late when Cape Wren stows during a pass. This morning Nkem registered `/events`, which answers 501 for now, and `/health`, which counts the streams `serveEvents` opens.

A server-sent event stream is an answer that stays open, `Content-Type: text/event-stream`, carrying events of a few lines and a blank line each. The server buffers what a handler writes, so each event must be flushed through the writer's `http.Flusher`, found by `http.NewResponseController(w)`. The request's context is done when the screen leaves.

## Your task

1. In `events.go`, write `writeEvent`: the view as one event, then `Flush`:

```text
event: board
id: 1
data: {"passes":[...],"weather":[...],...}

```

2. In `serveEvents`, set `Content-Type: text/event-stream` and `Cache-Control: no-store`, write `b.current()` with `writeEvent`, then wait until `r.Context()` is done.
