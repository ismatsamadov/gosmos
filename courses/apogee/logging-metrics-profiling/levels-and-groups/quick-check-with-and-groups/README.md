---
slug: quick-check-with-and-groups
title: Quick Check: With and Groups
type: quiz
position: 7
difficulty: 7
topics: logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:37Z
---
Thursday afternoon, before the new lines go to Lantern Atoll, Nkem writes five lines on the whiteboard and asks you three questions about them.

Here `o` holds the options, with `dropTime`:

```go
l := slog.New(slog.NewJSONHandler(os.Stdout, o))
u := l.With("program", "uplinkd")
g := slog.Group("pass", "sat_id", "brine-3")
u.Info("pass started", g)
l.Info("pass started", "site", "LAN")
```

## Your task

Answer Nkem's three questions: what the last line writes, how the text handler would write `g`, and where `program` lands in the fourth line's record.
