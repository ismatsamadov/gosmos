---
slug: quick-check-log-records
title: Quick Check: Log Records
type: quiz
position: 2
difficulty: 6
topics: json, logging
created: 2026-11-07T07:51:36Z
modified: 2026-11-11T17:33:38Z
---
Thursday. Before the new lines go to Lantern Atoll, Nkem writes four lines on the whiteboard and asks you three questions about them.

Here `o` holds the options, with `dropTime`:

```go
t := slog.New(slog.NewTextHandler(os.Stdout, o))
j := slog.New(slog.NewJSONHandler(os.Stdout, o))
t.Warn("carrier lost", "site", "LAN")
j.Warn("carrier lost", "site", "LAN", -1.5)
```

## Your task

Answer Nkem's three questions: where a station's code belongs in a record, what the last line writes after the site, and what changes between the lines the two handlers write for the same record.
