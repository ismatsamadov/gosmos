---
slug: quick-check-levels
title: Quick Check: Levels
type: quiz
position: 2
difficulty: 6
topics: logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-11T17:33:37Z
---
Tuesday, before the night shift gets the new setting, Ingrid has two questions: "What gets through. What happens when we change it mid-pass."

Levels are numbers, and a handler compares them:

```go
fmt.Println(int(slog.LevelDebug),
	int(slog.LevelInfo), int(slog.LevelWarn),
	int(slog.LevelError))
// -4 0 4 8
```

## Your task

Answer Ingrid's two questions: what a handler at WARN does with an INFO record, and what happens when the night shift changes a `LevelVar` while two passes are logging.
