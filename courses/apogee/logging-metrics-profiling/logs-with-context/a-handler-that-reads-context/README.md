---
slug: a-handler-that-reads-context
title: A Handler That Reads Context
type: code_tests
position: 1
difficulty: 7
topics: context, interfaces, logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
Tuesday. Nkem: "Nobody should have to remember the pass ID in a log call. The context already has it."

A `slog.Handler` is an interface of four methods, and a struct that embeds a handler gets all four from it; replacing `Handle` alone changes every record. `Handle` receives the context the call was given, `InfoContext(ctx, ...)`, and `Record.AddAttrs` adds attributes after the record's own:

```go
func (h siteHandler) Handle(ctx context.Context,
	r slog.Record) error {
	r.AddAttrs(slog.String("site", h.site))
	return h.Handler.Handle(ctx, r)
}
```

## Your task

In `passHandler`'s `Handle`, add `pass_id` from `passFrom` and `request_id` from `requestIDFrom`, each only when the context holds one, then hand the record to the handler inside. `requestIDFrom` comes from bookingd, the booking service you built: it reads the request IDs, such as `req-004211`, that bookingd puts in each request's context.

> [!NOTE]
> With the handler adding `pass_id`, `passLogger` keeps only the `pass` group from now on, so no line says its pass twice.
