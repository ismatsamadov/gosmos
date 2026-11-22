---
slug: quick-check-handlers
title: Quick Check: Handlers
type: quiz
position: 3
difficulty: 7
topics: context, interfaces, logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:38Z
---
Wednesday afternoon, before `passHandler` goes into every program, Mirela reads it: "Your handler adds the pass. Now show me where it doesn't."

In her five lines, `ctx` carries Cirrine-1's pass:

```go
h := passHandler{slog.NewJSONHandler(w, nil)}
l := slog.New(h)
l.Info("pass started")
l.WithGroup("radio").InfoContext(ctx,
	"carrier lost", "margin_db", -1.5)
```

## Your task

Answer Mirela's three questions: which context `Handle` gets for the third line, what a wrapper whose `WithAttrs` returns the inner handler's result loses, and where `pass_id` lands in the last record.
